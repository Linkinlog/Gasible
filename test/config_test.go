package test

import (
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
