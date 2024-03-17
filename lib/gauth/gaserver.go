package gauth

import (
	"fmt"
	"net"
	"net/http"
	"text/template"

	"github.com/rs/zerolog"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

var (
	host = "127.0.0.1"
	authEndpoint = "/codeResponse"
)

type GAToken struct {
	BindAddress		string
	AuthEndpoint	string
	Code 			string
	Config 			*oauth2.Config
}

type GAServer struct {
	State       string
	Code		chan string
	Listener    net.Listener
	Server      *http.Server
	GAToken		*GAToken
	Logger		*zerolog.Logger
}

type GAError struct {
	Name		string
	Description string
}

func newGAToken(clientID string, clientSecret string, port string) *GAToken {
	bindAddress := fmt.Sprintf("%s:%s", host, port)
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL: fmt.Sprintf("http://%s%s", bindAddress, authEndpoint),
		Endpoint:     google.Endpoint,
		Scopes:       []string{calendar.CalendarReadonlyScope},
	}
	return &GAToken{
		BindAddress: bindAddress,
		AuthEndpoint: authEndpoint,
		Config:      config,
	}
}

// NewGAServer makes the webserver for collecting auth
func NewGAServer(clientID string, clientSecret string, port string, logger *zerolog.Logger) *GAServer {
	return &GAServer{
		State:       "random-string",
		Logger:		 logger,
		GAToken:	 newGAToken(clientID, clientSecret, port),
		Code:		 make(chan string, 1),
	}
}

	// Reply with the response to the user and to the channel
func (s *GAServer) reply(w http.ResponseWriter, res *GAError) {
	var (
		status int
		responseTemplate string
	)
	if (res == nil) {
		status = http.StatusOK
	} else {
		status = http.StatusBadRequest
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/html")
	var t = template.Must(template.New("authResponse").Parse(responseTemplate))
	var err = t.Execute(w, res);
	if  err != nil {
		s.Logger.Debug().Msg(fmt.Sprintf("Could not execute template for web response: %s", err))
	}
}


func (s *GAServer) HandleAuth(w http.ResponseWriter, req *http.Request) {
	// Parse the form parameters and save them
	err := req.ParseForm()
	if err != nil {
		s.reply(w, &GAError{
			Name:        "Parse form error",
			Description: err.Error(),
		})
		return
	}

	// get code, error if empty
	var code = req.Form.Get("code")
	if code == "" {
		s.reply(w, &GAError{
			Name:        "Auth Error",
			Description: "No code returned by remote server",
		})
		return
	}

	// check state
	var state = req.Form.Get("state")
	if state != s.State {
		s.reply(w, &GAError{
			Name:        "Auth state doesn't match",
			Description: fmt.Sprintf("Expecting %q got %q", s.State, state),
		})
		return
	}

	// code OK
	s.reply(w, nil)
	s.Code <- req.FormValue("code")
}

// Init gets the internal web server ready to receive config details
func (gaServer *GAServer) Init() error {
	gaServer.Logger.Debug().Str("BindAddress", gaServer.GAToken.BindAddress).Msg("Starting auth server")
	var mux = http.NewServeMux()
	gaServer.Server = &http.Server{
		Addr:    gaServer.GAToken.BindAddress,
		Handler: mux,
	}
	gaServer.Server.SetKeepAlivesEnabled(false)

	mux.HandleFunc(gaServer.GAToken.AuthEndpoint, gaServer.HandleAuth)

	var err error
	gaServer.Listener, err = net.Listen("tcp", gaServer.GAToken.BindAddress)
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	return nil
}

// Serve the auth server, doesn't return
func (s *GAServer) Serve() error {
	var err = s.Server.Serve(s.Listener)
	return fmt.Errorf("closed auth server with error: %w", err)
}

// Stop the auth server by closing its socket
func (s *GAServer) Stop() {
	s.Logger.Debug().Msg("Closing auth server")
	s.Listener.Close()
	s.Server.Close()
}

func (gaServer *GAServer) Authorize() (*GAToken, error) {
	var err = gaServer.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to start auth webserver: %w", err)
	}

	go gaServer.Serve()
	defer gaServer.Stop()

	// Open the URL for the user to visit
	var authUrl = gaServer.GAToken.Config.AuthCodeURL(gaServer.State, oauth2.AccessTypeOffline)
	open.Start(authUrl)

	fmt.Printf("Waiting for gaServer.Code to be set\n")
	gaServer.GAToken.Code = <- gaServer.Code

	if	gaServer.GAToken.Code == "" {
		return nil, fmt.Errorf("failed to start auth webserver: %w", err)
	} 

	return gaServer.GAToken, nil
}

