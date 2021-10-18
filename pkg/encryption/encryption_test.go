package encryption_test

import (
	"testing"
	"bytes"

	"github.com/Iiqbal2000/let-us-lock/pkg/encryption"
)

func TestEncryptAES(t *testing.T) {
	plaintext := "The benchmark function must run the target code b.N times."
	key := "WytZ+j?0zIcUOQ@SsCBwW.ax_g*nQf(L"

	result, err := encryption.EncryptAes([]byte(plaintext), []byte(key))
	if err != nil {
		t.Fatal(err.Error())
	}
	
	if bytes.Equal(result, []byte(plaintext)) {
		t.Fatalf("Got error, actual value : %s", string(result))
	}
}