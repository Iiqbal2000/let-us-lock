package randstr_test

import (
	"testing"
	"os"

	"github.com/Iiqbal2000/let-us-lock/randstr"
)

func TestGenerateAndSave(t *testing.T) {
	expectLen := 100
	defer os.Remove("salt.txt")

	salt := randstr.GenerateAndSave(expectLen, "salt.txt")
	if len(salt) != expectLen {
		t.Errorf("Got: %d\nExpected: %d", len(salt), expectLen)
	}
}