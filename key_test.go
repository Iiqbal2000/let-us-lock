package main

import (
	"os"
	"path"
	"testing"
)

func removeDir(t *testing.T) {
	t.Helper()
	usr, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = os.RemoveAll(path.Join(usr, ".config", "let-us-lock"))
	if err != nil {
		t.Fatal(err.Error())
	}
}
func TestReadSaltFile(t *testing.T) {
	k, err := key{}.derive()
	if err != nil {
		t.Error(err.Error())
	}

	if len(k) == 0 {
		t.Fatal("the key should exist")
	}

	defer removeDir(t)
}
