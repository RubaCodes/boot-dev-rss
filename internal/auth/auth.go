package auth

import (
	"errors"
	"net/http"
	"strings"
)

// extract api_key from header
// example
// Authorization: ApiKey {#val}
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("authorization header not present")
	}
	vals := strings.Split(val, " ")
	if len(vals[1]) != 64 || vals[0] != "ApiKey" {
		return "", errors.New("invalid api key format")
	}
	return vals[1], nil
}
