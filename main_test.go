package main

import (
	"bytes"
	"log"
	"os"
	"path"
	"testing"

	"github.com/matryer/is"
)

var keytest = []byte("passphrasetesti\n")
var homeDir = getCfgPath()

func isFileExist(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		return false
	}
	return true
}

func deleteFiles(cfgPath string, ouputFile ...string) {
	if err := os.RemoveAll(cfgPath); err != nil {
		log.Fatal(err)
	}

	for _, f := range ouputFile {
		if err := os.Remove(f); err != nil {
			log.Fatal(err)
		}
	}
}

func TestEncrypt(t *testing.T) {
	cmd := []string{"main.go", "encrypt", "-f", "testdata/kitten.png", "-o", "testdata/cipherfile"}
	wantOutputFile := "testdata/cipherfile"
	wantSaltFile := path.Join(homeDir, cfgFile)

	is := is.New(t)

	r := bytes.NewBuffer(keytest)
	w := bytes.NewBuffer([]byte{})

	application := app{
		args:   cmd,
		input:  r,
		output: w,
	}
	err := application.runForTesting()

	is.NoErr(err)
	is.Equal(isFileExist(wantOutputFile), true)
	is.Equal(isFileExist(wantSaltFile), true)
	deleteFiles(homeDir, wantOutputFile)
}

func TestDecrypt(t *testing.T) {
	encryptCmd := []string{"main.go", "encrypt", "-f", "testdata/kitten.png", "-o", "testdata/cipherfile"}
	r := bytes.NewBuffer(keytest)
	w := bytes.NewBuffer([]byte{})

	application := app{
		args:   encryptCmd,
		input:  r,
		output: w,
	}

	err := application.runForTesting()
	if err != nil {
		t.Fatal(err.Error())
	}

	cmd := []string{"main.go", "decrypt", "-f", "testdata/cipherfile", "-o", "testdata/result.png"}
	wantOutputFile := "testdata/result.png"
	wantSaltFile := path.Join(homeDir, cfgFile)

	is := is.New(t)

	application = app{
		args:   cmd,
		input:  bytes.NewBuffer(keytest),
		output: w,
	}

	err = application.runForTesting()
	is.NoErr(err)
	is.Equal(isFileExist(wantOutputFile), true)
	is.Equal(isFileExist(wantSaltFile), true)
	deleteFiles(homeDir, wantOutputFile, "testdata/cipherfile")
}

func TestEncryptAndDecryptFailure(t *testing.T) {
	var cases = []struct {
		name string
		cmd  []string
		result string
	}{
		{
			name: "test with wrong command",
			cmd:  []string{"main.go", "ksiwn", "-f", "testdata/kitten.png", "-o", "testdata/cipherfile"},
			result: "testdata/cipherfile",
		},
		{
			name: "test without flags",
			cmd:  []string{"main.go", "encrypt"},
			result: "testdata/cipherfile",
		},
		{
			name: "test without command but flag exists",
			cmd:  []string{"main.go", "-f", "testdata/kitten.png", "-o", "testdata/cipherfile"},
			result: "testdata/cipherfile",
		},
		{
			name: "test without file flag",
			cmd:  []string{"main.go", "encrypt", "-o", "testdata/cipherfile"},
			result: "testdata/cipherfile",
		},
		{
			name: "test with two commands at a time",
			cmd:  []string{"main.go", "encrypt", "decrypt", "-f", "testdata/kitten.png", "-o", "testdata/cipherfile"},
			result: "testdata/cipherfile",
		},
	}

	for _, elem := range cases {
		t.Run(elem.name, func(t *testing.T) {
			r := bytes.NewBuffer(keytest)
			w := bytes.NewBuffer([]byte{})

			application := app{
				args:   elem.cmd,
				input:  r,
				output: w,
			}

			err := application.runForTesting()

			is := is.New(t)
			is.True(err != nil)
			is.True(isFileExist(elem.result) != true)
		})
	}
}
