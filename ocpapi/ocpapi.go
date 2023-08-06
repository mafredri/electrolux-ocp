package ocpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mafredri/electrolux-ocp/gigya"
	"golang.org/x/exp/slices"
)

const (
	APIURL = "https://api.ocp.electrolux.one"
)

type Client struct {
	config Config
	client *http.Client

	regionalBaseURL string
	authToken       AuthToken
}

type Config struct {
	APIKey       string
	Brand        string // Example: "electrolux"
	ClientID     string // Example: "ElxOneApp"
	ClientSecret string
	CountryCode  string // Example: "FI"
}

func New(config Config) (*Client, error) {
	if config.APIKey == "" {
		return nil, errors.New("missing APIKey")
	}
	if config.Brand == "" {
		return nil, errors.New("missing Brand")
	}
	if config.ClientID == "" {
		return nil, errors.New("missing ClientID")
	}
	if config.ClientSecret == "" {
		return nil, errors.New("missing ClientSecret")
	}
	if config.CountryCode == "" {
		return nil, errors.New("missing CountryCode")
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	c := &Client{
		client: httpClient,
		config: config,
	}
	httpClient.Transport = newClientTransport(http.DefaultTransport, config.APIKey)

	return c, nil
}

func (c *Client) SetAuthToken(token AuthToken) {
	c.authToken = token
}

func (c *Client) GetAuthToken() AuthToken {
	return c.authToken
}

// Login logs in to the API using the provided email and password.
func (c *Client) Login(ctx context.Context, email, password string) error {
	token, err := c.ClientToken(ctx)
	if err != nil {
		return fmt.Errorf("client token: %w", err)
	}

	ips, err := c.identityProviders(ctx, token, email)
	if err != nil {
		return fmt.Errorf("identity providers: %w", err)
	}

	switch len(ips) {
	case 0:
		return fmt.Errorf("no identity providers found")
	case 1:
	default:
		return fmt.Errorf("multiple identity providers found, only one is supported: found %d providers", len(ips))
	}
	ip := ips[0]
	c.regionalBaseURL = ip.HTTPRegionalBaseURL

	cs, err := c.countries(ctx, token)
	if err != nil {
		return fmt.Errorf("countries: %w", err)
	}

	if !slices.Contains(cs.CountryCodes(), c.config.CountryCode) {
		return fmt.Errorf("country code %q not found in available countries: %v", c.config.CountryCode, cs.CountryCodes())
	}

	gi := gigya.NewIdentity(gigya.Config{
		Domain: ip.Domain,
		APIKey: ip.APIKey,
	})
	idToken, err := gi.Login(ctx, email, password)
	if err != nil {
		return fmt.Errorf("gigya login: %w", err)
	}

	c.authToken, err = c.login(ctx, idToken)
	if err != nil {
		return fmt.Errorf("auth token: %w", err)
	}

	return nil
}

func (c *Client) ClientToken(ctx context.Context) (ClientToken, error) {
	body, err := json.Marshal(tokenRequest{
		GrantType:    "client_credentials",
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		Scope:        "",
	})
	if err != nil {
		return ClientToken{}, fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/one-account-authorization/api/v1/token", APIURL), bytes.NewReader(body))
	if err != nil {
		return ClientToken{}, fmt.Errorf("new request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return ClientToken{}, fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	var t ClientToken
	err = decodeResponse(resp, &t)
	if err != nil {
		return ClientToken{}, fmt.Errorf("decode response: %w", err)
	}

	return t, nil
}

func (c *Client) login(ctx context.Context, idToken string) (AuthToken, error) {
	body, err := json.Marshal(tokenRequest{
		GrantType: "urn:ietf:params:oauth:grant-type:token-exchange",
		ClientID:  c.config.ClientID,
		IDToken:   idToken,
		Scope:     "",
	})
	if err != nil {
		return AuthToken{}, fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/one-account-authorization/api/v1/token", c.regionalBaseURL), bytes.NewReader(body))
	if err != nil {
		return AuthToken{}, err
	}
	req.Header.Set("Origin-Country-Code", c.config.CountryCode)

	resp, err := c.client.Do(req)
	if err != nil {
		return AuthToken{}, fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	var t AuthToken
	err = decodeResponse(resp, &t)
	if err != nil {
		return AuthToken{}, err
	}

	return t, nil
}

func (c *Client) identityProviders(ctx context.Context, token ClientToken, email string) ([]identityProvider, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/one-account-user/api/v1/identity-providers?brand=%s&email=%s", APIURL, c.config.Brand, email), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token.Authorization())
	req.Header.Add("Context-Brand", c.config.Brand)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	var ips []identityProvider
	err = decodeResponse(resp, &ips)
	if err != nil {
		return nil, err
	}

	return ips, nil
}

func (c *Client) countries(ctx context.Context, token ClientToken) (countries, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/one-account-user/api/v1/countries", c.regionalBaseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token.Authorization())

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	var countries []country
	err = decodeResponse(resp, &countries)
	if err != nil {
		return nil, err
	}

	return countries, nil
}

func (c *Client) doAuthorized(req *http.Request) (*http.Response, error) {
	if !c.authToken.ExpiresAt.After(time.Now()) {
		// TODO(mafredri): Refresh.
		return nil, errors.New("auth token expired")
	}

	req.Header.Add("Authorization", c.authToken.Authorization())

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

// Appliances contains data from all appliances.
func (c *Client) Appliances(ctx context.Context, includeMetadata bool) ([]Appliance, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/appliance/api/v2/appliances?includeMetadata=%t", c.regionalBaseURL, includeMetadata), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.doAuthorized(req)
	if err != nil {
		return nil, fmt.Errorf("do authorized: %w", err)
	}
	defer resp.Body.Close()

	var appliances []Appliance
	err = decodeResponse(resp, &appliances)
	if err != nil {
		return nil, err
	}

	return appliances, nil
}

// AppliancesInfo contains information about the requested appliances.
func (c *Client) AppliancesInfo(ctx context.Context, applianceIDs ...string) ([]ApplianceInfo, error) {
	body, err := json.Marshal(map[string][]string{
		"applianceIds": applianceIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/appliance/api/v2/appliances/info", c.regionalBaseURL), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.doAuthorized(req)
	if err != nil {
		return nil, fmt.Errorf("do authorized: %w", err)
	}
	defer resp.Body.Close()

	var applianceInfo []ApplianceInfo
	err = decodeResponse(resp, &applianceInfo)
	if err != nil {
		return nil, err
	}

	return applianceInfo, nil
}

type clientTransport struct {
	rt     http.RoundTripper
	apiKey string
}

func newClientTransport(rt http.RoundTripper, apiKey string) http.RoundTripper {
	return &clientTransport{rt: rt, apiKey: apiKey}
}

func (ct *clientTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("x-api-key", ct.apiKey)
	req.Header.Add("User-Agent", "Ktor client")

	return ct.rt.RoundTrip(req)
}

func decodeResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code for %q: %d, body: %s", resp.Request.URL.Path, resp.StatusCode, string(b))
	}

	err := json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
