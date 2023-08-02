package internal

import (
	"fmt"
)

func ErrorAs(source string, err error) error {
	return fmt.Errorf("%w\nCaller: %s", err, source)
}
