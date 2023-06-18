package core_test

import (
	"testing"

	"github.com/Linkinlog/gasible/internal/core"
)

// Test that we can make a new module and register it
type MockModule struct {
	name string
	deps []string
}

func (mock MockModule) GetModuleDeps() []string                    { return mock.deps }
func (mock MockModule) ParseConfig(_ map[string]interface{}) error { return nil }
func (mock MockModule) Config() core.ModuleConfig                  { return core.ModuleConfig{} }
func (mock MockModule) Name() string                               { return mock.name }
func (mock MockModule) TearDown() error                            { return nil }
func (mock MockModule) Setup() error                               { return nil }
func (mock MockModule) Update() error                              { return nil }

type TestCase struct {
	TestName   string
	TestModule core.Module
}

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
	_, err = moduleRegistry.Get(testCase.TestModule.Name())
	if err != nil {
		t.Fatal(err)
	} else if _, ok := testCase.TestModule.(MockModule); !ok {
		// Test we can get the module we make
		t.Fatalf("Could not find module %v", testCase.TestModule)
	}
}

// TestTopologicalSort ensures that the modules are ran in topological order.
// Should always be [testModA testModB testModC testModD testModE lastModule].
func TestTopologicalSortInstallsInOrder(t *testing.T) {
	testModA := MockModule{
		name: "testModA",
	}
	testModB := MockModule{
		name: "testModB",
		deps: []string{"testModA"},
	}
	testModC := MockModule{
		name: "testModC",
		deps: []string{"testModB"},
	}
	testModD := MockModule{
		name: "testModD",
		deps: []string{"testModC"},
	}
	testModE := MockModule{
		name: "testModE",
		deps: []string{"testModD"},
	}
	lastModule := MockModule{
		name: "lastModule",
		deps: []string{"testModE"},
	}
	expectedResult := [6]string{"testModA", "testModB", "testModC", "testModD", "testModE", "lastModule"}
	modReg := core.NewModuleRegistry()
	modReg.Register(lastModule)
	modReg.Register(testModE)
	modReg.Register(testModB)
	modReg.Register(testModD)
	modReg.Register(testModA)
	modReg.Register(testModC)
	returnVal, err := modReg.TopologicallySortedModules()
	if err != nil {
		t.Fatal(err)
	}
	// Create an array of the names from the sorted modules
	var stringResult [6]string
	for i, module := range returnVal {
		stringResult[i] = module.Name()
	}
	if expectedResult != stringResult {
		t.Fatalf("modules out of order, expected %v, got %v", expectedResult, returnVal)
	} else {
		t.Log(returnVal)
	}
}
