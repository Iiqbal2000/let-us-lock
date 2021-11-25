package filesystem_test

import (
	"os"
	"testing"

	fs "github.com/Iiqbal2000/let-us-lock/filesystem"
)

func TestReadFile(t *testing.T) {
	testFile := "../kitten.png"
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatal(err)
	}
	expectSize := info.Size()
	
	result, err := fs.ReadFile(testFile)
	if err != nil {
		t.Fatal(err.Error())
	}

	if int64(len(result)) != expectSize {
		t.Fatalf("Got error, actual value: %d\nexpect: %d", result, expectSize)
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