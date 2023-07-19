package core_test

import (
	"testing"

	"github.com/Linkinlog/gasible/internal/core"
)

type MockModule struct {
	name   string
	deps   []string
	config core.ModuleConfig
}

func (mock MockModule) ParseConfig(_ map[string]interface{}) error { return nil }
func (mock MockModule) Config() core.ModuleConfig                  { return mock.config }
func (mock MockModule) GetName() string                            { return mock.name }
func (mock MockModule) TearDown() error                            { return nil }
func (mock MockModule) Setup() error                               { return nil }
func (mock MockModule) Update() error                              { return nil }

type TestCase struct {
	TestName   string
	TestModule core.Module
}

// Test that we can make a new module and register it
func TestRegisterAndGetNewModule(t *testing.T) {
	testCase := TestCase{
		TestName: "test",
		TestModule: MockModule{
			name: "TestModule",
		},
	}
	moduleRegistry := core.NewModuleRegistry()

	// Confirm there is nothing in the registry
	_, err := moduleRegistry.Get(testCase.TestName)
	if err != nil && err != core.ModuleNotFoundError {
		t.Fatalf("Expected error %s not found. Found %s", core.ModuleNotFoundError, err)
	}
	// Register
	moduleRegistry.Register(testCase.TestModule)
	// Confirm it is there now
	_, err = moduleRegistry.Get(testCase.TestModule.GetName())
	if err != nil {
		t.Fatal(err)
	} else if _, ok := testCase.TestModule.(MockModule); !ok {
		// Test we can get the module we make
		t.Fatalf("Could not find module %v", testCase.TestModule)
	}
}
