package bookworm

import (
	"testing"
)

func TestParse(t *testing.T) {
	if Parse() != "Huzzah" {
		t.Error("Expected huzzah")
	}
}
