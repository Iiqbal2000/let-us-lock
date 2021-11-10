package randstr_test

import (
	"testing"
	
	"github.com/Iiqbal2000/let-us-lock/randstr"
)

func TestGenerate(t *testing.T) {
	expectLen := 100

	salt := randstr.Generate(expectLen)
	if len(salt) != expectLen {
		t.Errorf("Got: %d\nExpected: %d", len(salt), expectLen)
	}
}