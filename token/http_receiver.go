package token

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HttpTokenReceiver struct {
	https bool
}

func NewHttpTokenReceiver() HttpTokenReceiver {
	return HttpTokenReceiver{
		https: false,
	}
}

func NewHttpsTokenReceiver() HttpTokenReceiver {
	return HttpTokenReceiver{
		https: true,
	}
}

func (r *HttpTokenReceiver) BuildConnectionString(server string, port uint, realm string, client string) (string, error) {
	var protocolStr string
	if r.https {
		protocolStr = "https"
	} else {
		protocolStr = "http"
	}
	return fmt.Sprintf("%s://%s:%d/realms/%s/protocol/openid-connect/token", protocolStr, server, port, realm), nil
}

func (r *HttpTokenReceiver) Get(connectionStr string, client string, password string, tokenReceiverChannel chan<- TokenReceiverPayload) {
	data := make(url.Values)

	data.Set("client_id", client)
	data.Set("grant_type", "client_credentials")
	data.Set("client_secret", password)

	req, err := http.NewRequest("POST", connectionStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		tokenReceiverChannel <- TokenReceiverPayloadError("Failed to create HTTP request: " + err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	clientHTTP := &http.Client{}
	resp, err := clientHTTP.Do(req)

	if err != nil {
		tokenReceiverChannel <- TokenReceiverPayloadError("HTTP request failed: " + err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tokenReceiverChannel <- TokenReceiverPayloadError("Error while reading HTTP response: " + err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		tokenReceiverChannel <- TokenReceiverPayloadError("HTTP request failed with status: " + resp.Status)
		return
	}

	if payload, err := stringToTokenReceiverPayload(body); err == nil {
		// Send result to channel
		tokenReceiverChannel <- payload
	} else {
		tokenReceiverChannel <- TokenReceiverPayloadError(err.Error())
	}
}
