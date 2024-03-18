/*
Copyright © 2024 Eugene Shtoka eshtoka@gmail.com

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/EugeneShtoka/gcalcli/lib/gauth"
	"github.com/EugeneShtoka/gcalcli/lib/levelwriter"
	"github.com/EugeneShtoka/gcalcli/lib/tokenrepo"
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
		getTokenFromWeb()
	},
}

// TODO: move port to optional params, and set default to 8080
// TODO: set command modes? to set, list and delete tokens
// TODO: add optional param for token name
// TODO: add optional param to set clientId and clientSecret from file
// TODO: add a help command with a list of available commands to usage
// TODO: add test cases

func getStringInput(paramName string, minLength int, maxLength int) string {
	var param string
	for {
		fmt.Printf("Enter %s: ", paramName)
		_, err := fmt.Scan(&param)

		if err == nil {
			if (len(param) < minLength) {
				if (len(param) == 0) {
					fmt.Printf("%s can't be empty\n", paramName)
				} else {
					fmt.Printf("%s must be at least %d characters long\n", paramName, minLength)
				}
				continue
			} else if (maxLength > 0 && len(param) > maxLength) {
				fmt.Printf("%s must be at most %d characters long\n", paramName, maxLength)
				continue
			}
			break
		} else {
			fmt.Println(err)
		}
	}
	return param
}

func getIntInput(paramName string, min *int, max *int) int {
	var param int
	for {
		fmt.Printf("Enter %s: ", paramName)
		_, err := fmt.Scan(&param)

		if err == nil {
			if (min != nil && param < *min) {
				fmt.Printf("%s must be greater than %d", paramName, *min)
				continue
			} else if (max != nil && param > *max) {
				fmt.Printf("%s must be greater than %d", paramName, *max)
				continue
			}
			break
		} else {
			fmt.Println(err)
		}
	}
	return param
}

func getTokenFromWeb() error { 
	var name = getStringInput("calendar name", 1, -1)
	var clientID = getStringInput("client id", 60, 100)
	var clientSecret = getStringInput("client secret", 30, 40)

	var minPort = 1024
	var maxPort = 65535
	var port = getIntInput("port", &minPort, &maxPort)

	var logger = levelwriter.NewLogger(zerolog.DebugLevel)
	return Authorize(name, clientID, clientSecret, fmt.Sprintf("%d", port), &logger)
}

func Authorize(name string, clientID string, clientSecret string, port string, logger *zerolog.Logger) error {
	server := gauth.NewGAServer(clientID, clientSecret, port, logger)
	token, err := server.Authorize()
	if err != nil {
		return fmt.Errorf("failed to authorize: %w", err)
	}
	
	fmt.Printf("Successfully authorized. Saving token named: %s to keyring.\n", name)
	err = tokenrepo.SaveToken(token, name)
	if err != nil {
		return fmt.Errorf("failed to save token %s to keyring: %w", name, err)
	}
	fmt.Printf("Token saved and ready for use\n")
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


