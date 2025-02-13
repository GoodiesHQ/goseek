package server

import (
	"slices"
	"sync"

	"github.com/goodieshq/goseek/utils"
)

type ApiKeyChecker interface {
	IsValidApiKey(apikey string) bool
	UpdateApiKeys(apikeys []string)
	AddApiKey(apikeys ...string)
	DelApiKey(apikeys ...string)
}

type ApiKeyCheckStatic struct {
	mu      sync.RWMutex
	apikeys []string
}

func NewApiKeyCheckStatic(apikeys []string) *ApiKeyCheckStatic {
	return &ApiKeyCheckStatic{
		apikeys: apikeys,
	}
}

func (apikeyCheck *ApiKeyCheckStatic) AddApiKey(apikeys ...string) {
	apikeyCheck.mu.Lock()
	defer apikeyCheck.mu.Unlock()
	apikeyCheck.apikeys = append(apikeyCheck.apikeys, apikeys...)
}

func (apikeyCheck *ApiKeyCheckStatic) DelApiKey(apikeys ...string) {
	apikeyCheck.mu.Lock()
	defer apikeyCheck.mu.Unlock()
	apikeyCheck.apikeys = utils.RemoveAll(apikeyCheck.apikeys, apikeys...)
}

func (apikeyCheck *ApiKeyCheckStatic) UpdateApiKeys(apikeys []string) {
	apikeyCheck.mu.Lock()
	defer apikeyCheck.mu.Unlock()
	apikeyCheck.apikeys = apikeys
}

func (apikeyCheck *ApiKeyCheckStatic) IsValidApiKey(apikey string) bool {
	apikeyCheck.mu.RLock()
	defer apikeyCheck.mu.RUnlock()
	return slices.Contains(apikeyCheck.apikeys, apikey)
}
