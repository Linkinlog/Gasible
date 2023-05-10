package modules

import "github.com/Linkinlog/gasible/internal/core"

type GenericPackageManager struct {
	Name: string,

}



func init() {
	core.ModuleRegistry.Register("GenericPackageManager", &GenericPackageManager{})
}

func (packageMan *GenericPackageManager) Setup() error {
	return nil
}

func (packageMan *GenericPackageManager) Update() error {
	return nil
}
