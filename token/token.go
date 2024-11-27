package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/okieoth/goptional"
)

type TokenContent struct {
	Token            string
	ExirationSeconds int64
	LastUpdated      goptional.Optional[time.Time]
	LastChecked      goptional.Optional[time.Time]
}

type TokenReceiver interface {
	Get(url string, client string, password string, tokenReceiverChannel <-chan string)
}

type Token struct {
	Url      string
	Client   string
	Password string
	Realm    string

	Content       goptional.Optional[TokenContent]
	TokenReceiver TokenReceiver
}

func (t *Token) Get() (string, error) {
	if c, isSet := t.Content.Get(); isSet {
		return c.Token, nil
	} else {
		return "", errors.New("Token not initialized")
	}
}

func NewTokenBuilder() TokenBuilder {
	var ret TokenBuilder
	return ret
}

type TokenBuilder struct {
	url      goptional.Optional[string]
	client   goptional.Optional[string]
	password goptional.Optional[string]
	realm    goptional.Optional[string]

	tokenReceiver goptional.Optional[TokenReceiver]
}

func (b *TokenBuilder) Url(v string) *TokenBuilder {
	b.url.Set(v)
	return b
}

func (b *TokenBuilder) Client(v string) *TokenBuilder {
	b.client.Set(v)
	return b
}

func (b *TokenBuilder) Password(v string) *TokenBuilder {
	b.password.Set(v)
	return b
}

func (b *TokenBuilder) Realm(v string) *TokenBuilder {
	b.realm.Set(v)
	return b
}

func (b *TokenBuilder) TokenReceiver(v TokenReceiver) *TokenBuilder {
	b.tokenReceiver.Set(v)
	return b
}

func (b *TokenBuilder) Build() (Token, error) {
	var ret Token
	if v, isSet := b.url.Get(); isSet {
		ret.Url = v
	} else {
		return ret, errors.New("url isn't set")
	}
	if v, isSet := b.client.Get(); isSet {
		ret.Client = v
	} else {
		return ret, errors.New("client isn't set")
	}
	if v, isSet := b.password.Get(); isSet {
		ret.Password = v
	} else {
		return ret, errors.New("password isn't set")
	}
	if v, isSet := b.realm.Get(); isSet {
		ret.Realm = v
	} else {
		return ret, errors.New("realm isn't set")
	}
	if v, isSet := b.tokenReceiver.Get(); isSet {
		ret.TokenReceiver = v
		tokenReceiverChannel := make(chan string)
		ret.TokenReceiver.Get(ret.Url, ret.Client, ret.Password, tokenReceiverChannel)
		for t := range tokenReceiverChannel {
			var tc TokenContent
			tc.Token = t
			ret.Content.Set(tc)
			// Start token refresh goroutine
		}
	} else {
		return ret, errors.New("token receiver isn't set")
	}

	return ret, nil
}

func Dummy() {
	fmt.Printf(":-) %v/n", time.Now())
}
