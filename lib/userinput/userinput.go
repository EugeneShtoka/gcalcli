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
package userinput

import (
	"fmt"
)

func GetStringInput(paramName string, minLength int, maxLength int) string {
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

func IsValidIntInput(paramName string, param int, min *int, max *int) bool {
	if (min != nil && param < *min) {
		fmt.Printf("%s must be greater than %d\n", paramName, *min)
		return false
	} else if (max != nil && param > *max) {
		fmt.Printf("%s must be greater than %d\n", paramName, *max)
		return false
	}

	return true
}

func GetIntInput(paramName string, min *int, max *int) int {
	var param int
	for {
		fmt.Printf("Enter %s: ", paramName)
		_, err := fmt.Scan(&param)

		if err == nil {
			if IsValidIntInput(paramName, param, min, max) {
				break
			}
			continue
		} else {
			fmt.Println(err)
		}
	}
	return param
}