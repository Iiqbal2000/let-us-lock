package aes_test

import (
	"testing"
	"bytes"

	"github.com/Iiqbal2000/let-us-lock/pkg/aes"
)

var aesTestCase = struct {
	plainText string
	key string
} {
	plainText: "The benchmark function must run the target code b.N times.",
	key: "WytZ+j?0zIcUOQ@SsCBwW.ax_g*nQf(L",
}

func TestEncrypt(t *testing.T) {
	result, err := aes.Encrypt([]byte(aesTestCase.plainText), []byte(aesTestCase.key))
	if err != nil {
		t.Fatal(err.Error())
	}
	
	if bytes.Equal(result, []byte(aesTestCase.plainText)) {
		t.Fatalf("Got error, actual value : %s", string(result))
	}
}

func TestDecrypt(t *testing.T) {
	
	cipherText, err := aes.Encrypt([]byte(aesTestCase.plainText), []byte(aesTestCase.key))
	if err != nil {
		t.Fatal(err.Error())
	}

	result, err := aes.Decrypt(cipherText, []byte(aesTestCase.key))
	if err != nil {
		t.Fatal(err.Error())
	}
	if !bytes.Equal(result, []byte(aesTestCase.plainText)) {
		t.Fatalf("Got error, actual value : %s", string(result))
	}
}