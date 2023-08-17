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

// State contains the current state of the client
// (e.g. for saving and restoring auth tokens).
type State struct {
	RegionalBaseURL string `json:"regionalBaseUrl"`
	ClientToken     Token  `json:"clientToken"`
	UserToken       Token  `json:"userToken"`
}

// Client is an Electrolux OCP API client.
type Client struct {
	config Config
	client *http.Client

	state State
}

type Config struct {
	APIURL       string
	APIKey       string
	Brand        string // Example: "electrolux"
	ClientID     string // Example: "ElxOneApp"
	ClientSecret string
	CountryCode  string // Example: "FI"

	State State // Optional initial state.
}

func New(config Config) (*Client, error) {
	if config.APIURL == "" {
		config.APIURL = APIURL
	}
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
		state:  config.State,
	}
	httpClient.Transport = newClientTransport(http.DefaultTransport, config.APIKey)

	return c, nil
}

// State returns the current state of the client.
func (c *Client) State() State {
	return c.state
}

type IdentityProvider struct {
	Domain                   string `json:"domain"` // "eu1.gigya.com"
	APIKey                   string `json:"apiKey"`
	Brand                    string `json:"brand"`                    // "electrolux"
	HTTPRegionalBaseURL      string `json:"httpRegionalBaseUrl"`      // "https://api.eu.ocp.electrolux.one"
	WebSocketRegionalBaseURL string `json:"webSocketRegionalBaseUrl"` // "wss://ws.eu.ocp.electrolux.one"
	DataCenter               string `json:"dataCenter"`               // "EU"
}

func (c *Client) IdentityProviders(ctx context.Context, email string) ([]IdentityProvider, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/one-account-user/api/v1/identity-providers?brand=%s&email=%s", c.config.APIURL, c.config.Brand, email), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Context-Brand", c.config.Brand)

	var ips []IdentityProvider
	err = c.doClientAuth(ctx, req, &ips)
	if err != nil {
		return nil, fmt.Errorf("do client auth: %w", err)
	}

	return ips, nil
}

type Country struct {
	Name           string `json:"name"`
	CountryCode    string `json:"countryCode"`
	LegalRegion    string `json:"legalRegion"`    // "APMEA", "EU (GDPR)", "LAM", "US"
	BusinessRegion string `json:"businessRegion"` // "BA-APMEA", "BA-EU", "BA-LATAM", "BA-NA"
	DataCenter     string `json:"dataCenter"`     // "AU", "US", "EU"
}

func countryCodes(cs ...Country) (codes []string) {
	for _, c := range cs {
		codes = append(codes, c.CountryCode)
	}
	return codes
}

func (c *Client) Countries(ctx context.Context) ([]Country, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/one-account-user/api/v1/countries", c.state.RegionalBaseURL), nil)
	if err != nil {
		return nil, err
	}

	var countries []Country
	err = c.doClientAuth(ctx, req, &countries)
	if err != nil {
		return nil, fmt.Errorf("do client auth: %w", err)
	}

	return countries, nil
}

// Login logs in to the API using the provided email and password.
func (c *Client) Login(ctx context.Context, email, password string) error {
	if c.state.RegionalBaseURL != "" && c.state.UserToken.RefreshToken != "" {
		// Assume a valid base URL and token has been provided.
		return nil
	}

	ips, err := c.IdentityProviders(ctx, email)
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
	c.state.RegionalBaseURL = ip.HTTPRegionalBaseURL

	if c.state.UserToken.RefreshToken != "" {
		// Assume a valid token has been provided.
		return nil
	}

	countries, err := c.Countries(ctx)
	if err != nil {
		return fmt.Errorf("countries: %w", err)
	}

	if codes := countryCodes(countries...); !slices.Contains(codes, c.config.CountryCode) {
		return fmt.Errorf("country code %q not found in available countries: %v", c.config.CountryCode, codes)
	}

	gi := gigya.NewIdentity(gigya.Config{
		Domain: ip.Domain,
		APIKey: ip.APIKey,
	})
	idToken, err := gi.Login(ctx, email, password)
	if err != nil {
		return fmt.Errorf("gigya login: %w", err)
	}

	c.state.UserToken, err = c.tokenExchange(ctx, idToken)
	if err != nil {
		return fmt.Errorf("auth token: %w", err)
	}

	return nil
}

// Appliances contains data from all appliances.
func (c *Client) Appliances(ctx context.Context, includeMetadata bool) ([]Appliance, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/appliance/api/v2/appliances?includeMetadata=%t", c.state.RegionalBaseURL, includeMetadata), nil)
	if err != nil {
		return nil, err
	}

	var appliances []Appliance
	err = c.doUserAuth(ctx, req, &appliances)
	if err != nil {
		return nil, fmt.Errorf("do user auth: %w", err)
	}

	return appliances, nil
}

// AppliancesInfo contains information about the requested appliances.
func (c *Client) AppliancesInfo(ctx context.Context, applianceIDs ...string) ([]ApplianceInfo, error) {
	body, err := json.Marshal(map[string][]string{"applianceIds": applianceIDs})
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/appliance/api/v2/appliances/info", c.state.RegionalBaseURL), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	var applianceInfo []ApplianceInfo
	err = c.doUserAuth(ctx, req, &applianceInfo)
	if err != nil {
		return nil, fmt.Errorf("do user auth: %w", err)
	}

	return applianceInfo, nil
}

// Token is an access token generated via client or user authentication.
type Token struct {
	AccessToken  string    `json:"accessToken"`            // JWT.
	ExpiresIn    int       `json:"expiresIn"`              // Seconds.
	ExpiresAt    time.Time `json:"expiresAt"`              // Set by UnmarshalJSON (now + expires in).
	TokenType    string    `json:"tokenType"`              // "Bearer"
	RefreshToken string    `json:"refreshToken,omitempty"` // Set for login.
	Scope        string    `json:"scope"`                  // Example: "", "email offline_access eluxiot:*:*:*"
}

// Authorization returns the Authorization header value for the token.
func (t Token) Authorization() string {
	return fmt.Sprintf("%s %s", t.TokenType, t.AccessToken)
}

// UnmarshalJSON implements json.Unmarshaler and assigns ExpiresAt.
func (t *Token) UnmarshalJSON(b []byte) error {
	now := time.Now()

	type token Token
	var tmp token
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*t = Token(tmp)
	if t.ExpiresAt.IsZero() {
		t.ExpiresAt = now.Add(time.Duration(t.ExpiresIn) * time.Second)
	}
	return nil
}

type tokenRequest struct {
	GrantType    string `json:"grantType"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret,omitempty"` // Set for ClientToken.
	IDToken      string `json:"idToken,omitempty"`      // Set for login.
	RefreshToken string `json:"refreshToken,omitempty"` // Set for refresh.
	Scope        string `json:"scope"`
}

func (c *Client) clientCredentials(ctx context.Context) (Token, error) {
	body, err := json.Marshal(tokenRequest{
		GrantType:    "client_credentials",
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		Scope:        "",
	})
	if err != nil {
		return Token{}, fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/one-account-authorization/api/v1/token", c.config.APIURL), bytes.NewReader(body))
	if err != nil {
		return Token{}, fmt.Errorf("new request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	var t Token
	err = c.do(ctx, req, &t)
	if err != nil {
		return Token{}, fmt.Errorf("do: %w", err)
	}

	return t, nil
}

func (c *Client) tokenExchange(ctx context.Context, idToken string) (Token, error) {
	return c.token(ctx, tokenRequest{
		GrantType: "urn:ietf:params:oauth:grant-type:token-exchange",
		ClientID:  c.config.ClientID,
		IDToken:   idToken,
		Scope:     "",
	})
}

func (c *Client) refreshToken(ctx context.Context, token Token) (Token, error) {
	if token.RefreshToken == "" {
		return Token{}, errors.New("refresh token is missing")
	}

	return c.token(ctx, tokenRequest{
		GrantType:    "refresh_token",
		ClientID:     c.config.ClientID,
		RefreshToken: token.RefreshToken,
		Scope:        "",
	})
}

func (c *Client) token(ctx context.Context, tr tokenRequest) (Token, error) {
	body, err := json.Marshal(tr)
	if err != nil {
		return Token{}, fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/one-account-authorization/api/v1/token", c.state.RegionalBaseURL), bytes.NewReader(body))
	if err != nil {
		return Token{}, err
	}
	req.Header.Add("Authorization", "Bearer ")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Origin-Country-Code", c.config.CountryCode)

	var t Token
	err = c.do(ctx, req, &t)
	if err != nil {
		return Token{}, fmt.Errorf("do: %w", err)
	}

	return t, nil
}

func (c *Client) doClientAuth(ctx context.Context, req *http.Request, v any) error {
	if c.state.ClientToken.AccessToken == "" || time.Now().After(c.state.ClientToken.ExpiresAt) {
		var err error
		c.state.ClientToken, err = c.clientCredentials(ctx)
		if err != nil {
			return fmt.Errorf("client token: %w", err)
		}
	}

	req.Header.Add("Authorization", c.state.ClientToken.Authorization())

	return c.do(ctx, req, v)
}

func (c *Client) doUserAuth(ctx context.Context, req *http.Request, v any) error {
	if c.state.UserToken.AccessToken == "" {
		return errors.New("please login before using this endpoint")
	}
	if time.Now().After(c.state.UserToken.ExpiresAt) {
		authToken, err := c.refreshToken(ctx, c.state.UserToken)
		if err != nil {
			return fmt.Errorf("auth token expired: refresh failed: %w", err)
		}
		c.state.UserToken = authToken
	}

	req.Header.Add("Authorization", c.state.UserToken.Authorization())

	return c.do(ctx, req, v)
}

func (c *Client) do(ctx context.Context, req *http.Request, v any) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("http client do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code for %q: %d, body: %s", resp.Request.URL.Path, resp.StatusCode, string(b))
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
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
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Charset", "UTF-8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")

	// Go's http transport automatically requests gzip.
	// req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	return ct.rt.RoundTrip(req)
}
