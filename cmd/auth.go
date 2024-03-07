/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("auth called")
		config := &oauth2.Config{
			ClientID:     "1096710297832-tmmj2s9tfj5v27hbsmbq7b7i7tecisht.apps.googleusercontent.com",
			ClientSecret: "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",
			RedirectURL: "http://localhost:8080",
			Endpoint:     google.Endpoint,
			Scopes:       []string{calendar.CalendarReadonlyScope},
		}
		getTokenFromWeb(config)
	},
}

var codeChannel = make(chan string, 1)
var oauthStateString = "pseudo-random"
var authConfig oauth2.Config

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
    state := r.FormValue("state")
	authCode := r.FormValue("code")
	fmt.Printf("Code: %s, State: %s, oauthStateString: %s\n", authCode, state, oauthStateString)
    if state != oauthStateString {
		log.Fatal("Invalid oauth state.")
    }

	fmt.Printf("authCode: %s\n", authCode)
	tok, err := authConfig.Exchange(context.TODO(), authCode)
	if err != nil {
			log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	codeChannel <- authCode

	fmt.Printf("authCode: %v\n", tok)

}

func getTokenFromWeb(config *oauth2.Config) error { 
	authConfig = *config
	authURL := config.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	open.Start(authURL)
	// fmt.Printf("Go to the following link in your browser then type the "+
	//         "authorization code: \n%v\n", authURL)

	http.HandleFunc("/", handleGoogleCallback)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	server := http.Server{
		Addr: ":8080",
	}

	go server.Serve(listener)
	fmt.Printf("Served & Listening on %s\n", listener.Addr())
	result := <-codeChannel
	fmt.Printf("Result: %s\n", result)
	
	close(codeChannel)
	listener.Close()
	server.Close()

	return nil
}



func init() {
	rootCmd.AddCommand(authCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
