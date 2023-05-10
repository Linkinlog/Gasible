package yamlParser

import (
	"os"

	"github.com/Linkinlog/gasible/internal/core"
	"gopkg.in/yaml.v3"
)

// CreateDefaults will generate a YAML file
// using the defaults we outline.
func CreateDefaults(file string) error {
	if file == "" {
		file = "default.yml"
	}
	Conf := core.NewConfigWithDefaults()
	d, err := yaml.Marshal(&Conf)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, d, 0644)
	if err != nil {
		return err
	}

	return nil
}
