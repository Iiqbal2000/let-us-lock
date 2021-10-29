package filesystem_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	fs "github.com/Iiqbal2000/let-us-lock/filesystem"
)

func TestReadFile(t *testing.T) {
	expectation := []byte("temporary file's content")
	tmpfile, err := ioutil.TempFile(".", "example")
	if err != nil {
		t.Fatal(err.Error())
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(expectation); err != nil {
		t.Fatal(err.Error())
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err.Error())
	}

	result, err := fs.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err.Error())
	}
	if !bytes.Equal(expectation, result) {
		t.Fatalf("Got error, actual value: %s", string(result))
	}
}

func TestWriteFile(t *testing.T) {
	fileNameExpect := "example.txt"
	contentExpect := []byte("new file's content")

	result, err := fs.WriteFile(contentExpect, fileNameExpect)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer os.Remove(result)
	
	if result != fileNameExpect {
		t.Fatalf("Got error, actual value: %s", result)
	}
}