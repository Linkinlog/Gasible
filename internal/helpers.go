package internal

import (
	"fmt"
)

// ErrorAs is a formatter for errors, so we can see what function had what error.
func ErrorAs(source string, err error) error {
	return fmt.Errorf("%w\nCaller: %s", err, source)
}
