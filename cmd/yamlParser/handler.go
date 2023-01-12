package yamlParser

import (
	"os"

	"github.com/Linkinlog/gasible/internal/models"
	"gopkg.in/yaml.v3"
)

// CreateDefaults will generate a YAML file
// using the defaults we outline.
func CreateDefaults() error {
	Conf := models.Config{}.Default()
	d, err := yaml.Marshal(&Conf)
	if err != nil {
		return err
	}

	err = os.WriteFile("default.yml", d, 0644)
	if err != nil {
		return err
	}

	return nil
}
