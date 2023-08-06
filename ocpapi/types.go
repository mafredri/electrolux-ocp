package ocpapi

import (
	"encoding/json"
	"fmt"
	"time"
)

type tokenRequest struct {
	GrantType    string `json:"grantType"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret,omitempty"` // Set for apiToken.
	IDToken      string `json:"idToken,omitempty"`      // Set for authToken.
	Scope        string `json:"scope"`
}

type ClientToken struct {
	AccessToken string    `json:"accessToken"`
	ExpiresIn   int       `json:"expiresIn"`
	ExpiresAt   time.Time `json:"-"`
	TokenType   string    `json:"tokenType"` // "Bearer"
	Scope       string    `json:"scope"`     // Example: "", "email offline_access eluxiot:*:*:*"
}

func (t ClientToken) Authorization() string {
	return fmt.Sprintf("%s %s", t.TokenType, t.AccessToken)
}

type AuthToken struct {
	ClientToken
	RefreshToken string `json:"refreshToken"`
}

func (t *ClientToken) UnmarshalJSON(b []byte) error {
	type token ClientToken
	var tmp token
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*t = ClientToken(tmp)
	t.ExpiresAt = time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)
	return nil
}

type countries []country

func (c countries) CountryCodes() (codes []string) {
	for _, country := range c {
		codes = append(codes, country.CountryCode)
	}
	return codes
}

type country struct {
	Name           string `json:"name"`
	CountryCode    string `json:"countryCode"`
	LegalRegion    string `json:"legalRegion"`    // "APMEA", "EU (GDPR)", "LAM", "US"
	BusinessRegion string `json:"businessRegion"` // "BA-APMEA", "BA-EU", "BA-LATAM", "BA-NA"
	DataCenter     string `json:"dataCenter"`     // "AU", "US", "EU"
}

type identityProvider struct {
	Domain                   string `json:"domain"` // "eu1.gigya.com"
	APIKey                   string `json:"apiKey"`
	Brand                    string `json:"brand"`                    // "electrolux"
	HTTPRegionalBaseURL      string `json:"httpRegionalBaseUrl"`      // "https://api.eu.ocp.electrolux.one"
	WebSocketRegionalBaseURL string `json:"webSocketRegionalBaseUrl"` // "wss://ws.eu.ocp.electrolux.one"
	DataCenter               string `json:"dataCenter"`               // "EU"
}
