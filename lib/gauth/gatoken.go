package gauth

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type GAToken struct {
	Name 			string	
	BindAddress		string
	AuthEndpoint	string
	Code 			string
	Config 			*oauth2.Config
	Logger 			*zerolog.Logger
}

var (
	localhost = "http://localhost"
	authEndpoint = "/codeResponse"
)

func newGAToken(name string, clientID string, clientSecret string, port string, logger *zerolog.Logger) *GAToken {
	bindAddress := fmt.Sprintf("%s:%s", localhost, port)
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL: fmt.Sprintf("%s%s", bindAddress, authEndpoint),
		Endpoint:     google.Endpoint,
		Scopes:       []string{calendar.CalendarReadonlyScope},
	}
	return &GAToken{
		Name:        name,
		BindAddress: bindAddress,
		AuthEndpoint: authEndpoint,
		Config:      config,
		Logger:		 logger,
	}
}

func (token *GAToken) get() (*oauth2.Token, error) {
	if (token.Code == "") {
		return nil, fmt.Errorf("Token was not authenticated")
	}
	return token.Config.Exchange(context.TODO(), token.Code)
}