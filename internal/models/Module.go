// Module file
//
// This file contains everything directly related to handling a Module model.
package models

import (
	"errors"
)

const MODULE_NOT_FOUND_ERROR = "No Module Found."

type Module interface {
	Setup() error
	Update() error
}

type ModuleRegistry struct {
	modules map[string]Module
}

func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		make(map[string]Module),
	}
}

func (mr *ModuleRegistry) Get(mod string) (Module, error) {
	found := mr.modules[mod]
	if found != nil {
		return found, nil
	}
	return nil, errors.New(MODULE_NOT_FOUND_ERROR)
}

func (mr *ModuleRegistry) Register(name string, mod Module) {
	mr.modules[name] = mod
}
