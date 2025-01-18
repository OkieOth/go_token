package main

import (
	"fmt"

	"github.com/okieoth/gotoken/token"

	"time"
)

func main() {
	tokenReceiver := token.NewHttpTokenReceiver()
	if token, err := token.NewTokenBuilder().
		Server("localhost").
		Port(8080).
		Realm("test-realm").
		Client("test-client").
		Password("test-client999").
		TokenReceiver(&tokenReceiver).
		Build(); err != nil {
		panic(fmt.Sprintf("Error while create token object: %v", err))
	} else {
		// do something useful with the token
		getAndPrintToken(&token)
		for {
			fmt.Println("sleep for 30s ... |-) ... ")
			time.Sleep(time.Second * 30)
			getAndPrintToken(&token)
		}
	}
}

func getAndPrintToken(token *token.Token) {
	tokenContent, err := token.Get()
	if err != nil {
		panic(fmt.Sprintf("Error while get token content: %v", err))
	}
	fmt.Println("Token Content: ", tokenContent)
	fmt.Println()
}
