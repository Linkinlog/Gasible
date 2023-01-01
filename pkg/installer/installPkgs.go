package installer

import (
	"errors"

	"github.com/Linkinlog/gasible/pkg/osHandler"
)

func Installer(s osHandler.System) error {
    return errors.New("WIP")
    // use s.pkgManager, attempt to install the packages
}
