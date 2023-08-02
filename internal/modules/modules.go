package modules

import "github.com/Linkinlog/gasible/internal/app"

type module app.Module

var ToBeRegistered []module

var ToBeInstalled map[packageManager][]string = make(map[packageManager][]string)
