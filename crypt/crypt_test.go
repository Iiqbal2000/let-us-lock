package crypt_test

import (
	"testing"
	"bytes"

	"github.com/Iiqbal2000/let-us-lock/crypt"
)

var aesTestCase = struct {
	plainText string
	key string
} {
	plainText: "The benchmark function must run the target code b.N times.",
	key: "WytZ+j?0zIcUOQ@SsCBwW.ax_g*nQf(L",
}

func TestEncrypt(t *testing.T) {
	result, err := crypt.Encrypt([]byte(aesTestCase.plainText), []byte(aesTestCase.key))
	if err != nil {
		t.Fatal(err.Error())
	}
	
	if bytes.Equal(result, []byte(aesTestCase.plainText)) {
		t.Fatalf("Got error, actual value : %s", string(result))
	}
}

func TestDecrypt(t *testing.T) {
	
	cipherText, err := crypt.Encrypt([]byte(aesTestCase.plainText), []byte(aesTestCase.key))
	if err != nil {
		t.Fatal(err.Error())
	}

	result, err := crypt.Decrypt(cipherText, []byte(aesTestCase.key))
	if err != nil {
		t.Fatal(err.Error())
	}
	if !bytes.Equal(result, []byte(aesTestCase.plainText)) {
		t.Fatalf("Got error, actual value : %s", string(result))
	}
}