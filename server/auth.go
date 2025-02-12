package server

import "slices"

type AuthChecker interface {
	CheckBasic(username, password string) bool
	CheckApiKey(apikey string) bool
}

type AuthCheck struct {
	credentials map[string]string
	apikeys     []string
}

func NewAuthCheck(credentials map[string]string, apikeys []string) *AuthCheck {
	return &AuthCheck{
		credentials: credentials,
		apikeys:     apikeys,
	}
}

func (authcheck *AuthCheck) CheckBasic(username, password string) bool {
	pw, found := authcheck.credentials[username]
	return found && password == pw
}

func (authcheck *AuthCheck) CheckApiKey(apikey string) bool {
	return slices.Contains(authcheck.apikeys, apikey)
}
