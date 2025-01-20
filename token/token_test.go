package token_test

import (
	"testing"
	"time"

	"github.com/okieoth/gotoken/token"
)

func TestGetToken_IT(t *testing.T) {
	tokenReceiver := token.NewHttpTokenReceiver()
	if token, err := token.NewToken(
		token.Realm("test-realm"),
		token.Client("test-client"),
		token.Password("test-client999"),
		token.Receiver(&tokenReceiver)); err != nil {
		t.Errorf("Error while create token object: %v", err)
	} else {
		// do something useful with the token
		tokenContent, err := token.Get()
		if err != nil {
			t.Errorf("Error while get token content: %v", err)
			return
		}
		if len(tokenContent) == 0 {
			t.Errorf("No token content: %v", err)
			return
		}
	}

}

func TestTokenUpdate_IT(t *testing.T) {
	tokenReceiver := token.NewHttpTokenReceiver()
	if token, err := token.NewToken(
		token.Realm("test-realm"),
		token.Client("test-client"),
		token.Password("test-client999"),
		token.Receiver(&tokenReceiver)); err != nil {
		t.Errorf("Error while create token object: %v", err)
	} else {
		// do something useful with the token
		tokenContent, err := token.Get()
		if err != nil {
			t.Errorf("Error while get token content: %v", err)
			return
		}
		if len(tokenContent) == 0 {
			t.Errorf("No token content: %v", err)
			return
		}
	}
	time.Sleep(time.Minute * time.Duration(5))
}
