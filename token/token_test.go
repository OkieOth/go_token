package token_test

import (
	"testing"

	"github.com/okieoth/gotoken/token"
)

func TestGetToken_IT(t *testing.T) {
	tokenReceiver := token.NewHttpTokenReceiver()
	if token, err := token.NewTokenBuilder().
		Server("localhost").
		Port(8080).
		Realm("test-realm").
		Client("test-client").
		Password("test-client999").
		TokenReceiver(&tokenReceiver).
		Build(); err != nil {
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
