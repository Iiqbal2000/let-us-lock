package main

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	expectLen := 1000
	
	salt := GenerateSalt(expectLen)
	if len(salt) != expectLen {
		t.Errorf("Got: %d\nExpected: %d", len(salt), expectLen)
	}
}