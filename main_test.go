package main

import (
	"bytes"
	"log"
	"os"
	"path"
	"testing"

	"github.com/matryer/is"
)

var passphraseTest = []byte("passphrasetesti\n")
var cfgDir = path.Join(getHomeDir(), cfgDirDefault)

func isFileExist(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		return false
	}
	return true
}

func deleteFiles(fileNames ...string) {
	if err := os.RemoveAll(cfgDir); err != nil {
		log.Fatal(err)
	}

	for _, f := range fileNames {
		err := os.Remove(f)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestEncryptDecrypt(t *testing.T) {
	t.Cleanup(func() {
		deleteFiles("testdata/result.png", "testdata/cipherfile")
	})

	t.Run("Encrypt", func(t *testing.T) {
		cmd := []string{"main.go", "encrypt", "-f", "testdata/kitten.png", "-o", "testdata/cipherfile"}
		wantOutputFile := "testdata/cipherfile"
		wantSaltFile := path.Join(cfgDir, saltFileName)

		is := is.New(t)

		r := bytes.NewBuffer(passphraseTest)
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
	})

	t.Run("Decrypt", func(t *testing.T) {
		cmd := []string{"main.go", "decrypt", "-f", "testdata/cipherfile", "-o", "testdata/result.png"}
		wantOutputFile := "testdata/result.png"
		wantSaltFile := path.Join(cfgDir, saltFileName)

		is := is.New(t)

		application := app{
			args:   cmd,
			input:  bytes.NewBuffer(passphraseTest),
			output: bytes.NewBuffer([]byte{}),
		}

		err := application.runForTesting()
		is.NoErr(err)

		is.Equal(isFileExist(wantOutputFile), true)
		is.Equal(isFileExist(wantSaltFile), true)
	})
}

func TestEncryptDecryptFailure(t *testing.T) {
	t.Cleanup(func() {
		deleteFiles()
	})

	var cases = []struct {
		name   string
		cmd    []string
		result string
	}{
		{
			name:   "test with wrong command",
			cmd:    []string{"main.go", "ksiwn", "-f", "testdata/kitten.png", "-o", "testdata/cipherfile"},
			result: "testdata/cipherfile",
		},
		{
			name:   "test without flags",
			cmd:    []string{"main.go", "encrypt"},
			result: "testdata/cipherfile",
		},
		{
			name:   "test without command but flag exists",
			cmd:    []string{"main.go", "-f", "testdata/kitten.png", "-o", "testdata/cipherfile"},
			result: "testdata/cipherfile",
		},
		{
			name:   "test without file flag",
			cmd:    []string{"main.go", "encrypt", "-o", "testdata/cipherfile"},
			result: "testdata/cipherfile",
		},
		{
			name:   "test with two commands at a time",
			cmd:    []string{"main.go", "encrypt", "decrypt", "-f", "testdata/kitten.png", "-o", "testdata/cipherfile"},
			result: "testdata/cipherfile",
		},
	}

	for _, elem := range cases {
		t.Run(elem.name, func(t *testing.T) {
			is := is.New(t)

			r := bytes.NewBuffer(passphraseTest)
			w := bytes.NewBuffer([]byte{})

			application := app{
				args:   elem.cmd,
				input:  r,
				output: w,
			}

			err := application.runForTesting()
			is.True(err != nil)

			is.True(isFileExist(elem.result) != true)
		})
	}
}
