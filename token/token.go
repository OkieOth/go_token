package token

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type TokenContent struct {
	mutex            sync.RWMutex
	token            string
	ExirationSeconds uint64
	LastUpdated      *time.Time
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

type TokenOpts struct {
	server        string
	port          uint
	client        string
	password      string
	realm         string
	tokenReceiver *TokenReceiver
}

type TokenOptsFunc func(*TokenOpts)

func defaultTokenOpts() TokenOpts {
	return TokenOpts{
		server: "localhost",
		port:   8080,
	}
}

func Server(server string) TokenOptsFunc {
	return func(opts *TokenOpts) {
		opts.server = server
	}
}

func Port(port uint) TokenOptsFunc {
	return func(opts *TokenOpts) {
		opts.port = port
	}
}

func Client(client string) TokenOptsFunc {
	return func(opts *TokenOpts) {
		opts.client = client
	}
}

func Password(pwd string) TokenOptsFunc {
	return func(opts *TokenOpts) {
		opts.password = pwd
	}
}

func Realm(realm string) TokenOptsFunc {
	return func(opts *TokenOpts) {
		opts.realm = realm
	}
}

func Receiver(receiver TokenReceiver) TokenOptsFunc {
	return func(opts *TokenOpts) {
		opts.tokenReceiver = &receiver
	}
}

type Token struct {
	TokenOpts

	connectionStr string
	Content       *TokenContent
}

func NewToken(opts ...TokenOptsFunc) (*Token, error) {
	o := defaultTokenOpts()
	for _, fn := range opts {
		fn(&o)
	}
	ret := Token{
		connectionStr: "",
		Content:       nil,
		TokenOpts:     o,
	}
	errorMsg := ""
	if ret.client == "" {
		errorMsg += "no client config given,"
	}
	if ret.password == "" {
		errorMsg += " no password given,"
	}
	if ret.realm == "" {
		errorMsg += " no realm given,"
	}
	if ret.tokenReceiver == nil {
		errorMsg += " tokenReceiver not initialized"
	}
	if errorMsg != "" {
		return nil, fmt.Errorf("Missing items for creating a token: %s", errorMsg)
	} else {
		if err := ret.initializeToken(); err != nil {
			return nil, fmt.Errorf("error while initialize: %v", err)
		} else {
			return &ret, nil
		}
	}
}

func (t *Token) initializeToken() error {
	tokenReceiverChan := make(chan TokenReceiverPayload)
	connectionString, err := (*t.tokenReceiver).BuildConnectionString(t.server, t.port, t.realm, t.client)
	if err != nil {
		return fmt.Errorf("error while building connection string: %v", err)
	}
	t.connectionStr = connectionString
	go (*t.tokenReceiver).Get(connectionString, t.client, t.password, tokenReceiverChan)
	timeout := 10 * time.Second
	select {
	case payload := <-tokenReceiverChan:
		if payload.HasError {
			return errors.New(payload.Error)
		} else {
			t.InitContent(payload)
			go refreshToken(t)
		}
	case <-time.After(timeout):
		return errors.New("timeout while receiving the first token")
	}
	return nil
}

func (t *Token) Get() (string, error) {
	if t.Content != nil {
		return t.Content.GetToken(), nil
	} else {
		return "", errors.New("Token not initialized")
	}
}

func refreshToken(t *Token) {
	const EXPIRATION_OFFSET = 5
	const MAX_RETRY_SECS = 60
	currentRetrySecs := 1
	for {
		expirationSecs := t.Content.ExirationSeconds
		if expirationSecs > EXPIRATION_OFFSET {
			time.Sleep(time.Second * time.Duration(t.Content.ExirationSeconds-EXPIRATION_OFFSET))
		} else {
			time.Sleep(time.Second * time.Duration(currentRetrySecs))
			if currentRetrySecs < MAX_RETRY_SECS {
				currentRetrySecs = currentRetrySecs * 2
			}
		}
		tokenReceiverChan := make(chan TokenReceiverPayload)
		go (*t.tokenReceiver).Get(t.connectionStr, t.client, t.password, tokenReceiverChan)
		timeout := 10 * time.Second
		select {
		case payload := <-tokenReceiverChan:
			if payload.HasError {
				// TODO Logging
				log.Printf("error while retrieving token: %v", payload.Error)
				t.Content.ExirationSeconds = 0
			} else {
				log.Println("retrieved new token")
				if t.Content != nil {
					t.Content.SetToken(payload.TokenStr)
					lastUpdated := time.Now()
					t.Content.LastUpdated = &lastUpdated
					currentRetrySecs = 1
				} else {
					t.InitContent(payload)
				}
				log.Println("updated token")
			}
		case <-time.After(timeout):
			log.Println("timeout while trying to refresh token")
			t.Content.ExirationSeconds = 0
		}
	}
}

func (t *Token) InitContent(payload TokenReceiverPayload) {
	var content TokenContent
	content.ExirationSeconds = payload.ExpirationSeconds
	content.SetToken(payload.TokenStr)
	t.Content = &content
}
