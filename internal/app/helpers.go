package app

import (
	"fmt"
)

func ErrorAs(source string, err error) error {
	return fmt.Errorf("%w occured in %s", err, source)
}
