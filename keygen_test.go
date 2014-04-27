package main

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Errorf("TestGenerateKey error: %s key returned: %s", err, key)
	}
}
