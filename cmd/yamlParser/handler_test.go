package yamlParser_test

import (
	"testing"

	"github.com/Linkinlog/gasible/cmd/yamlParser"
)

// TODO fix bad tests
func TestCreateDefaults(t *testing.T) {
	t.Skip("Test needs fixed, creates file and doesnt clean up")
	t.Run("TestCreateDefaults", func(t *testing.T) {
		if err := yamlParser.CreateDefaults(); err != nil {
			t.Fatalf("Failed to create YAML defaults, got: %s", err.Error())
		}
	})

}
