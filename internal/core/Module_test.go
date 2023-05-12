package core_test

import (
	"testing"

	"github.com/Linkinlog/gasible/internal/core"
)

// Test that we can make a new module and register it
type MockModule struct{}

func (mock MockModule) Setup() error {
	return nil
}

func (mock MockModule) Update() error {
	return nil
}

type TestCase struct {
	TestName   string
	TestModule core.Module
}

func TestRegisterAndGetNewModule(t *testing.T) {
	testCase := TestCase{
		"test",
		MockModule{},
	}
	moduleRegistry := core.NewModuleRegistry()

	// Confirm there is nothing in the registry
	_, err := moduleRegistry.Get(testCase.TestName)
	if err.Error() != core.MODULE_NOT_FOUND_ERROR {
		t.Fatalf("Expected error %s not found. Found %s", core.MODULE_NOT_FOUND_ERROR, err)
	}
	// Register
	moduleRegistry.Register(testCase.TestName, testCase.TestModule)
	// Confirm it is there now
	module, err := moduleRegistry.Get(testCase.TestName)
	if err != nil {
		t.Fatal(err)
	} else if module != testCase.TestModule {
		// Test we can get the module we make
		t.Fatalf("Could not find module %v", testCase.TestModule)
	}
}
