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
package levelwriter

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
    // Using different time formats to make it clear which logger is being used
    debugOut = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123} // time format: "Mon, 02 Jan 2006 15:04:05 MST"
    errorOut = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339} // time format: "2006-01-02T15:04:05Z07:00"
)

// logOut implements zerolog.LevelWriter
type logOut struct{
    minLevel zerolog.Level
}

// Write should not be 


func (l logOut) Write(p []byte) (n int, err error) {
    return os.Stdout.Write(p)
}

// WriteLevel write to the appropriate output
func (l logOut) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
    if level >= zerolog.ErrorLevel {
        return errorOut.Write(p)
    } else if level >= l.minLevel {
        return debugOut.Write(p)
    }
    return 0, nil;
}

func NewLogger(minLevel zerolog.Level) zerolog.Logger {
    return zerolog.New(logOut{minLevel}).With().Timestamp().Logger()
}