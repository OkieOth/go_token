package token

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type TokenContent struct {
	mutex            sync.RWMutex
	token            string
	ExirationSeconds uint64
	LastUpdated      *time.Time
	LastChecked      *time.Time
}

func (c *TokenContent) GetToken() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.token
}

func (c *TokenContent) SetToken(t string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.token = t
}

type TokenReceiverPayload struct {
	TokenStr          string
	ExpirationSeconds uint64
	HasError          bool
	Error             string
}

func TokenReceiverPayloadError(msg string) TokenReceiverPayload {
	return TokenReceiverPayload{
		HasError: true,
		Error:    msg,
	}
}

type TokenReceiver interface {
	BuildConnectionString(server string, port uint, realm string, client string) (string, error)
	Get(connectionString string, client string, password string, tokenReceiverChannel chan<- TokenReceiverPayload)
}

type Token struct {
	Server   string
	Port     uint
	Client   string
	Password string
	Realm    string

	Content       *TokenContent
	TokenReceiver TokenReceiver
}

func (t *Token) Get() (string, error) {
	if t.Content != nil {
		return t.Content.GetToken(), nil
	} else {
		return "", errors.New("Token not initialized")
	}
}

func (t *Token) InitContent(payload TokenReceiverPayload) {
	var content TokenContent
	content.ExirationSeconds = payload.ExpirationSeconds
	content.SetToken(payload.TokenStr)
	t.Content = &content
	go func() {

	}()
	// TODO - initialize Token object
	// start go routine to refresh the token
}

func NewTokenBuilder() *TokenBuilder {
	var ret TokenBuilder
	return &ret
}

type TokenBuilder struct {
	server   *string
	port     *uint
	client   *string
	password *string
	realm    *string

	tokenReceiver *TokenReceiver
}

func (b *TokenBuilder) Server(v string) *TokenBuilder {
	b.server = &v
	return b
}

func (b *TokenBuilder) Port(v uint) *TokenBuilder {
	b.port = &v
	return b
}

func (b *TokenBuilder) Client(v string) *TokenBuilder {
	b.client = &v
	return b
}

func (b *TokenBuilder) Password(v string) *TokenBuilder {
	b.password = &v
	return b
}

func (b *TokenBuilder) Realm(v string) *TokenBuilder {
	b.realm = &v
	return b
}

func (b *TokenBuilder) TokenReceiver(v TokenReceiver) *TokenBuilder {
	b.tokenReceiver = &v
	return b
}

func (b *TokenBuilder) Build() (Token, error) {
	var ret Token
	if b.server != nil {
		ret.Server = *b.server
	} else {
		return ret, errors.New("server isn't set")
	}
	if b.port != nil {
		ret.Port = *b.port
	} else {
		return ret, errors.New("port isn't set")
	}
	if b.client != nil {
		ret.Client = *b.client
	} else {
		return ret, errors.New("client isn't set")
	}
	if b.password != nil {
		ret.Password = *b.password
	} else {
		return ret, errors.New("password isn't set")
	}
	if b.realm != nil {
		ret.Realm = *b.realm
	} else {
		return ret, errors.New("realm isn't set")
	}
	if b.tokenReceiver != nil {
		ret.TokenReceiver = *b.tokenReceiver
		tokenReceiverChan := make(chan TokenReceiverPayload)
		connectionString, err := ret.TokenReceiver.BuildConnectionString(ret.Server, ret.Port, ret.Realm, ret.Client)
		if err != nil {
			return ret, fmt.Errorf("error while building connection string: %v", err)
		}
		go ret.TokenReceiver.Get(connectionString, ret.Client, ret.Password, tokenReceiverChan)
		timeout := 10 * time.Second
		select {
		case payload := <-tokenReceiverChan:
			if payload.HasError {
				return ret, errors.New(payload.Error)
			} else {
				ret.InitContent(payload)
			}
		case <-time.After(timeout):
			return ret, errors.New("timeout while receiving the first token")
		}
	} else {
		return ret, errors.New("token receiver isn't set")
	}

	return ret, nil
}
