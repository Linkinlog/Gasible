package yamlParser

import (
	"os"

	"github.com/Linkinlog/gasible/internal/core"
	"gopkg.in/yaml.v3"
)

// CreateDefaults will generate a YAML file
// using the defaults we outline.
func CreateDefaults() error {
	Conf := core.NewConfigWithDefaults()
	d, err := yaml.Marshal(&Conf)
	if err != nil {
		return err
	}

	err = os.WriteFile("config.yml", d, 0644)
	if err != nil {
		return err
	}

	return nil
}
