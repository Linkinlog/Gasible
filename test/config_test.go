package test

import (
	"os"
	"testing"

	"github.com/Linkinlog/gasible/internal/models"
)

// These following tests test that our FillFromFile method works as intended.
// We want it to use defaults or panic returning an error
// informing the user of said error. (Likely a missing config)
func TestConfigFillFromFileWithNoConfig(t *testing.T) {
	cf := models.Config{}.Default()
	if err := cf.FillFromFile(); err == nil {
		t.Fatalf("Failed TestConfigFillFromFileWithNoConfig: %v", err)
	}
}

// Test the FillFromFile method whenever there is no config file
func TestConfigFillFromFileWithConfig(t *testing.T) {
	cf := models.Config{}.Default()
	file, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	cf.GlobalOpts.FilePath = file.Name()
	if err := cf.FillFromFile(); err != nil {
		t.Fatal(err)
	}

}
