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
package tokenrepo

import (
	"encoding/json"
	"fmt"

	"github.com/EugeneShtoka/gcalcli/lib/gauth"
	"github.com/zalando/go-keyring"
)

var serviceName = "gcalcli"

func LoadToken(calendarName string) (*gauth.GAToken, error) {
	var tokenString, err = keyring.Get(serviceName, calendarName)
	if err != nil {
        return nil, fmt.Errorf("failed to load token: %w", err)
    }

	var target gauth.GAToken
    err = json.Unmarshal([]byte(tokenString), &target)
    if err != nil {
        return nil, fmt.Errorf("error deserializing object: %w" + err.Error())
    }

    return &target, nil
}

func SaveToken(gaToken *gauth.GAToken, calendarName string) error {
    var jsonData, err = json.Marshal(gaToken)
    if err != nil {
        return fmt.Errorf("failed to serialize token to JSON: %w", err)
    }

    var jsonStr = string(jsonData) 
	err = keyring.Set(serviceName, calendarName, jsonStr)
	if err != nil {
        return fmt.Errorf("failed to save token: %w", err)
    }
	
	return nil
}