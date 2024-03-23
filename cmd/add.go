/*
Copyright © 2024 Eugene Shtoka <eshtoka@gmail.com>

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

	"github.com/EugeneShtoka/gcalcli/lib/gauth"
	"github.com/EugeneShtoka/gcalcli/lib/levelwriter"
	"github.com/EugeneShtoka/gcalcli/lib/tokenrepo"
	"github.com/EugeneShtoka/gcalcli/lib/userinput"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var defaultPort int = 58080

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
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

func getTokenFromWeb() error { 
	// var name = userinput.GetStringInput("calendar name", 1, -1)
	// var clientID = userinput.GetStringInput("client id", 60, 100)
	// var clientSecret = userinput.GetStringInput("client secret", 30, 40)

	var logger = levelwriter.NewLogger(zerolog.DebugLevel)
	// return authorize(name, clientID, clientSecret, fmt.Sprintf("%d", getPort()), &logger)
	return authorize("name", "clientID", "clientSecret", fmt.Sprintf("%d", getPort()), &logger)
}

func getPort() int {
	var minPort = 1024
	var maxPort = 65535
	var portName = "port"
	var port = viper.GetInt("port")
	if  (userinput.IsValidIntInput(portName, port, &minPort, &maxPort)) {
		return port
	}
	return userinput.GetIntInput(portName, &minPort, &maxPort)
}

func authorize(name string, clientID string, clientSecret string, port string, logger *zerolog.Logger) error {
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
	authCmd.AddCommand(addCmd)

	viper.SetDefault("port", defaultPort)
	addCmd.Flags().IntP("port", "p", defaultPort, "port number for gAuth code response")
	viper.BindPFlag("port", addCmd.Flags().Lookup("port"))
}
