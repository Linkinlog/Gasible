package modules

import "github.com/Linkinlog/gasible/internal/app"

// module is defined here so we can use it easier within our modules.
type module app.Module

// ToBeRegistered is a slice of modules that we will register upon application start.
var ToBeRegistered []module

// ToBeInstalled is for modules to set which package manager needs to install which dependencies.
var ToBeInstalled = make(map[packageManager][]string)
