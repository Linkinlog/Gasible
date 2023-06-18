package core_test

import (
	"testing"

	"github.com/Linkinlog/gasible/internal/core"
)

// Test that we can make a new module and register it
type MockModule struct {
	deps []string
}

func (mock MockModule) GetDeps() []string {
	return mock.deps
}

func (mock MockModule) ParseConfig(_ map[string]interface{}) error {
	return nil
}

func (mock MockModule) Config() core.ModuleConfig {
	return core.ModuleConfig{}
}

func (mock MockModule) TearDown() error {
	return nil
}

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
		TestName:   "test",
		TestModule: MockModule{},
	}
	moduleRegistry := core.NewModuleRegistry()

	// Confirm there is nothing in the registry
	_, err := moduleRegistry.Get(testCase.TestName)
	if err != nil && err != core.ModuleNotFoundError {
		t.Fatalf("Expected error %s not found. Found %s", core.ModuleNotFoundError, err)
	}
	// Register
	moduleRegistry.Register(testCase.TestName, testCase.TestModule)
	// Confirm it is there now
	_, err = moduleRegistry.Get(testCase.TestName)
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
	expectedResult := [6]string{"testModA", "testModB", "testModC", "testModD", "testModE", "lastModule"}
	testModA := MockModule{}
	testModB := MockModule{
		deps: []string{"testModA"},
	}
	testModC := MockModule{
		deps: []string{"testModB"},
	}
	testModD := MockModule{
		deps: []string{"testModC"},
	}
	testModE := MockModule{
		deps: []string{"testModD"},
	}
	lastModule := MockModule{
		deps: []string{"testModE"},
	}
	modReg := core.NewModuleRegistry()
	modReg.Register("lastModule", lastModule)
	modReg.Register("testModA", testModA)
	modReg.Register("testModB", testModB)
	modReg.Register("testModC", testModC)
	modReg.Register("testModD", testModD)
	modReg.Register("testModE", testModE)
	returnVal, err := modReg.TopologicallySortedModuleDeps()
	if err != nil {
		t.Fatal(err)
	} else if expectedResult != ([6]string)(returnVal) {
		t.Fatalf("modules out of order, expected %v, got %v", expectedResult, returnVal)
	} else {
		t.Log(returnVal)
	}
}
