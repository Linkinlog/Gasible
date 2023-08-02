package app_test

import (
	"testing"

	"github.com/Linkinlog/gasible/internal/app"
)

type MockModule struct {
	name   string
	config app.ModuleConfig
}

func (mock MockModule) ParseConfig(_ map[string]interface{}) error { return nil }
func (mock MockModule) Config() app.ModuleConfig                   { return mock.config }
func (mock MockModule) GetName() string                            { return mock.name }
func (mock MockModule) TearDown() error                            { return nil }
func (mock MockModule) Setup() error                               { return nil }
func (mock MockModule) Update() error                              { return nil }
func (mock MockModule) SetApp(_ *app.App)                          {}

type TestCase struct {
	TestName   string
	TestModule app.Module
}

// Test that we can make a new module and register it
func TestRegisterAndGetNewModule(t *testing.T) {
	testCase := TestCase{
		TestName: "test",
		TestModule: MockModule{
			name: "TestModule",
		},
	}
	application := app.New()

	// Confirm there is nothing in the registry
	_, err := application.ModuleRegistry.GetModule(testCase.TestName)
	if err != nil && err != app.ModuleNotFoundError {
		t.Fatalf("Expected error %s not found. Found %s", app.ModuleNotFoundError, err)
	}
	// Register
	application.ModuleRegistry.Register(testCase.TestModule)
	// Confirm it is there now
	_, err = application.ModuleRegistry.GetModule(testCase.TestModule.GetName())
	if err != nil {
		t.Fatal(err)
	} else if _, ok := testCase.TestModule.(MockModule); !ok {
		// Test we can get the module we make
		t.Fatalf("Could not find module %v", testCase.TestModule)
	}
}
