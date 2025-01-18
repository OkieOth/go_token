package token

import (
	"encoding/json"
	"fmt"
)

func stringToTokenReceiverPayload(body []byte) (TokenReceiverPayload, error) {
	var ret TokenReceiverPayload
	type tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   uint64 `json:"expires_in"`
	}
	var tr tokenResponse
	if err := json.Unmarshal(body, &tr); err == nil {
		ret.ExpirationSeconds = tr.ExpiresIn
		ret.TokenStr = tr.AccessToken
		return ret, nil
	} else {
		return ret, fmt.Errorf("error while unmarshal received token response: %v\nbody: %s", err, string(body))
	}
}
