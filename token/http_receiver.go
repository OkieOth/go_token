package token

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/okieoth/goptional"
)

type HttpTokenReceiver struct {
}

func (r *HttpTokenReceiver) Get(urlStr string, client string, password string, tokenReceiverChannel chan<- TokenReceiverPayload) {
	data := make(url.Values)

	data.Set("client_id", client)
	data.Set("grant_type", "client_credentials")
	data.Set("client_secret", password)

	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		e := goptional.NewOptionalValue[string]("Failed to create HTTP request: " + err.Error())
		tokenReceiverChannel <- TokenReceiverPayload{
			Error: e,
		}
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	clientHTTP := &http.Client{}
	resp, err := clientHTTP.Do(req)

	if err != nil {
		e := goptional.NewOptionalValue[string]("HTTP request failed: " + err.Error())
		tokenReceiverChannel <- TokenReceiverPayload{
			Error: e,
		}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		e := goptional.NewOptionalValue[string]("Error while reading HTTP response: " + err.Error())
		tokenReceiverChannel <- TokenReceiverPayload{
			Error: e,
		}
		return
	}

	if resp.StatusCode != http.StatusOK {
		e := goptional.NewOptionalValue[string](fmt.Sprintf("HTTP request failed with status: %s", resp.Status))
		tokenReceiverChannel <- TokenReceiverPayload{
			Error: e,
		}
		return
	}

	if payload, err := stringToTokenReceiverPayload(body); err == nil {
		// Send result to channel
		tokenReceiverChannel <- payload
	} else {
		var e goptional.Optional[string]
		e.Set(err.Error())
		tokenReceiverChannel <- TokenReceiverPayload{
			Error: e,
		}
	}
}
