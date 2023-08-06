package gigya

type response struct {
	CallID       string `json:"callId"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	APIVersion   int    `json:"apiVersion"`
	StatusCode   int    `json:"statusCode"`
	StatusReason string `json:"statusReason"`
	Time         string `json:"time"`
}

type loginResponse struct {
	response
	RegisteredTimestamp        int         `json:"registeredTimestamp"`
	Uid                        string      `json:"UID"`
	UIDSignature               string      `json:"UIDSignature"`
	SignatureTimestamp         string      `json:"signatureTimestamp"`
	Created                    string      `json:"created"`
	CreatedTimestamp           int         `json:"createdTimestamp"`
	IsActive                   bool        `json:"isActive"`
	IsRegistered               bool        `json:"isRegistered"`
	IsVerified                 bool        `json:"isVerified"`
	LastLogin                  string      `json:"lastLogin"`
	LastLoginTimestamp         int         `json:"lastLoginTimestamp"`
	LastUpdated                string      `json:"lastUpdated"`
	LastUpdatedTimestamp       int         `json:"lastUpdatedTimestamp"`
	LoginProvider              string      `json:"loginProvider"`
	OldestDataUpdated          string      `json:"oldestDataUpdated"`
	OldestDataUpdatedTimestamp int         `json:"oldestDataUpdatedTimestamp"`
	Profile                    profile     `json:"profile"`
	Registered                 string      `json:"registered"`
	SocialProviders            string      `json:"socialProviders"`
	Verified                   string      `json:"verified"`
	VerifiedTimestamp          int         `json:"verifiedTimestamp"`
	NewUser                    bool        `json:"newUser"`
	SessionInfo                sessionInfo `json:"sessionInfo"`
}

type profile struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Email     string `json:"email"`
	Zip       string `json:"zip"`
}

type sessionInfo struct {
	SessionToken  string `json:"sessionToken"`
	SessionSecret string `json:"sessionSecret"`
	ExpiresIn     string `json:"expires_in"`
}

type jwtResponse struct {
	response
	IDToken string `json:"id_token"`
}
