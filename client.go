package goverdrive

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/mrjones/oauth"
)

//Client is the overdrive API Client object
type Client struct {
	Consumer    *oauth.Consumer
	user        *oauth.AccessToken
	client      *http.Client
	consumerKey string
}

func GetOAuthURL(clientID, clientSecret, libraryAccountId, code, redirectURI string) (token *oauth2.Token) {
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{libraryAccountId},
		Endpoint: oauth2.Endpoint{
			AuthURL:  AuthURL,
			TokenURL: TokenURL,
		},
	}

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	return tok
}

//NewClient is the constructor with only the Consumer key and secret
func NewClient(key string, secret string) *Client {
	c := Client{
		consumerKey: key,
	}
	c.SetConsumer(key, secret)
	return &c
}

// Constructor with Consumer key/secret and user token/secret
func NewClientWithToken(consumerKey, consumerSecret, token, tokenSecret string) *Client {
	c := NewClient(consumerKey, consumerSecret)
	c.SetToken(token, tokenSecret)
	return c
}

// Set Consumer credentials, invalidates any previously cached client
func (c *Client) SetConsumer(consumerKey string, consumerSecret string) {
	c.Consumer = oauth.NewConsumer(consumerKey, consumerSecret,
		oauth.ServiceProvider{})
	c.client = nil
}

// Set user credentials, invalidates any previously cached client
func (c *Client) SetToken(token string, secret string) {
	c.user = &oauth.AccessToken{
		AdditionalData: nil,
		Secret:         secret,
		Token:          token,
	}
	c.client = nil
}

// Retrieve the underlying HTTP client
func (c *Client) GetHttpClient() (*http.Client, error) {
	if c.Consumer == nil {
		return nil, errors.New("Consumer credentials are not set")

	}
	if c.user == nil {
		c.SetToken("", "")
	}
	if c.client == nil {
		c.client, _ = c.Consumer.MakeHttpClient(c.user)
		c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	return c.client, nil
}
