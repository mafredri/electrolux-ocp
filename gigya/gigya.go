// Package gigya implements the Gigya/SAP login for use by ocpapi.
package gigya

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	Domain string
	APIKey string
}

type Identity struct {
	config Config
	client *http.Client
}

func NewIdentity(config Config) *Identity {
	return &Identity{
		config: config,
		client: http.DefaultClient,
	}
}

func (v *Identity) Login(ctx context.Context, user, password string) (jwtToken string, err error) {
	l, err := v.login(ctx, user, password)
	if err != nil {
		return "", err
	}

	jt, err := v.jwtToken(ctx, l.Uid, l.SessionInfo.SessionToken, l.SessionInfo.SessionSecret)
	if err != nil {
		return "", err
	}

	return jt, nil
}

func (v *Identity) login(ctx context.Context, user, password string) (loginResponse, error) {
	form := url.Values{
		"apikey":          []string{v.config.APIKey},
		"format":          []string{"json"},
		"httpStatusCodes": []string{"false"},
		"loginID":         []string{user},
		"password":        []string{password},
		"targetEnv":       []string{"mobile"},
	}
	uri := fmt.Sprintf("https://accounts.%s/accounts.login", v.config.Domain)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, strings.NewReader(form.Encode()))
	if err != nil {
		return loginResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.client.Do(req)
	if err != nil {
		return loginResponse{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var res loginResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return loginResponse{}, fmt.Errorf("decode response: %w", err)
	}
	if res.ErrorCode > 0 {
		return loginResponse{}, fmt.Errorf("login: %s", res.ErrorMessage)
	}

	return res, nil
}

func (v *Identity) jwtToken(ctx context.Context, uid, sessionToken, sessionSecret string) (string, error) {
	form := url.Values{
		"apikey":          []string{v.config.APIKey},
		"fields":          []string{"country"},
		"format":          []string{"json"},
		"httpStatusCodes": []string{"false"},
		"targetUID":       []string{uid},
		"oauth_token":     []string{sessionToken},
		"secret":          []string{sessionSecret},
		"targetEnv":       []string{"mobile"},
	}
	uri := fmt.Sprintf("https://accounts.%s/accounts.getJWT", v.config.Domain)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var res jwtResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if res.ErrorCode > 0 {
		return "", fmt.Errorf("jwtToken: %s", res.ErrorMessage)
	}

	return res.IDToken, err
}
