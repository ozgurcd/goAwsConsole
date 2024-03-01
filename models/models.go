package models

type AwsCredentials struct {
	SessionId    string `json:"sessionId"`
	SessionKey   string `json:"sessionKey"`
	SessionToken string `json:"sessionToken"`
}

type AwsFederationResponse struct {
	SigninToken string `json:"SigninToken"`
}

type RuntimeConfig struct {
	Profile     string
	RoleName    string
	Duration    int32
	Region      string
	Browser     string
	SeparateWin bool
	ProfileDir  string
}
