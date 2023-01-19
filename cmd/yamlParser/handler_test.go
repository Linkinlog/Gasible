package yamlParser_test

import (
	"testing"

	"github.com/Linkinlog/gasible/cmd/yamlParser"
)

func TestCreateDefaults(t *testing.T) {
    type testCase struct{
        // Filepath for testing
        f string
    }
    testfile := t.TempDir() + "test.yml"
    t1 := testCase{
        f: testfile,
    }
    t.Run("TestCreateDefaults", func(t *testing.T) {
        if err := yamlParser.CreateDefaults(t1.f); err != t1.ExpectedResponse {
            t.Fatalf("Failed to create YAML defaults, got: %s", err.Error())
        }
    })

}
