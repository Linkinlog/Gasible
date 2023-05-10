// Module file
//
// This file contains everything directly related to handling a Module model.
package models

import (
	"errors"
)

type Module interface {
	Setup() error
	Update() error
}

