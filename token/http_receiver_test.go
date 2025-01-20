package token_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/okieoth/gotoken/token"
)

func TestGetOk_IT(t *testing.T) {
	test := []struct {
		name       string
		port       uint
		realm      string
		client     string
		password   string
		expiredSec uint64
		isError    bool
	}{
		{"ok", 8080, "test-realm", "test-client", "test-client999", 23, false},
		{"wrong port", 9000, "test-realm", "test-client", "test-client999", 23, true},
		{"wrong realm", 8080, "XXXX", "test-client", "test-client999", 23, true},
		{"wrong client", 8080, "test-realm", "XXXX", "test-client999", 23, true},
		{"wrong password", 8080, "test-realm", "test-client", "XXXX", 23, true},
	}

	server, set := os.LookupEnv("KEYCLOAK_HOST")
	if !set {
		server = "localhost"
	}

	timeout := 2 * time.Second
	for _, testCase := range test {
		resultChan := make(chan token.TokenReceiverPayload)
		var receiver token.HttpTokenReceiver
		urlStr := fmt.Sprintf("http://%s:%d/realms/%s/protocol/openid-connect/token", server, testCase.port, testCase.realm)
		go receiver.Get(urlStr, testCase.client, testCase.password, resultChan)
		select {
		case payload := <-resultChan:
			if payload.HasError != testCase.isError {
				t.Error("Error is set")
			}
			if (!testCase.isError) && (payload.TokenStr == "") {
				t.Error("No token string found")
			}
			if (!testCase.isError) && (payload.ExpirationSeconds != testCase.expiredSec) {
				t.Errorf("Wrong expiration seconds: %d", payload.ExpirationSeconds)
			}
		case <-time.After(timeout):
			t.Error("Timeout while receiving the first token")
		}
	}
}
