package main

import (
	"bytes"
	"log"
	"os"
	"testing"
)

var keytest = []byte("passphrasetesti\n")

var testSuccessCase = []struct{
	name string
	cmd []string
	wantFile []string
} {
		{
			name: "test encrypt",
			cmd: []string{"main.go", "encrypt", "-f", "testdata/kitten.png", "-o", "testdata/cipherfile"},
			wantFile: []string{"salt.txt", "testdata/cipherfile"},
		},
		{	
			name: "test decrypt",
			cmd: []string{"main.go", "decrypt", "-f", "testdata/cipherfile", "-o", "testdata/result.png"},
			wantFile: []string{"salt.txt", "testdata/result.png"},
		},
}

var testFailureCase = []struct{
	name string
	cmd []string
	file []string
} {
	{
		name: "test with wrong command",
		cmd: []string{"main.go", "ksiwn", "-f", "testdata/kitten.png", "-o", "testdata/cipherfile"},
		file: []string{"salt.txt", "testdata/cipherfile"},
	},
	{
		name: "test without flags",
		cmd: []string{"main.go", "encrypt"},
		file: []string{"salt.txt", "testdata/cipherfile"},
	},
	{
		name: "test without command but flag is exist",
		cmd: []string{"main.go", "-f", "testdata/kitten.png", "-o", "testdata/cipherfile"},
		file: []string{"salt.txt", "testdata/cipherfile"},
	},
	{
		name: "test without file flag",
		cmd: []string{"main.go", "encrypt", "-o", "testdata/cipherfile"},
		file: []string{"salt.txt", "testdata/cipherfile"},
	},
	{
		name: "test with two commands at a time",
		cmd: []string{"main.go", "encrypt", "decrypt", "-f", "testdata/kitten.png", "-o", "testdata/cipherfile"},
		file: []string{"salt.txt", "testdata/cipherfile"},
	},
}

func isFileExist(filename string) error {
	if _, err := os.Stat(filename); err != nil {
		return err
	}
	return nil
}

func deleteFiles() {
	if err := os.Remove("salt.txt"); err != nil {
		log.Fatal(err)
	}
	if err := os.Remove("testdata/cipherfile"); err != nil {
		log.Fatal(err)
	}
	if err := os.Remove("testdata/result.png"); err != nil {
		log.Fatal(err)
	}
}

func TestEncryptAndDecrypt(t *testing.T) {
	var stdin bytes.Buffer
	
	for _, val := range testSuccessCase {
		stdin.Write(keytest)
		
		if val.cmd[1] == "encrypt"{
			t.Run(val.name, func(t *testing.T) {
				if err := run(val.cmd, &stdin); err != nil {
					t.Fatalf("failure, something wrong\nerror: %v", err)
				} else {
					if err := isFileExist(val.wantFile[0]); err != nil {
						t.Fatalf("failure: want: %v file is not found\nactual : %v", val.wantFile[0], err)	
					}

					if err := isFileExist(val.wantFile[1]); err != nil {
						t.Fatalf("failure: want: %v file is not found\nactual : %v", val.wantFile[1], err)	
					}
				}
			})
		} else if val.cmd[1] == "decrypt" {
			t.Run(val.name, func(t *testing.T) {
				if err := run(val.cmd, &stdin); err != nil {
					t.Fatalf("failure, something wrong\nerror: %v", err)
				} else {
					if err := isFileExist(val.wantFile[1]); err != nil {
						t.Fatalf("failure: want: %v file is not found\nactual : %v", val.wantFile[1], err)	
					}
				}

				func() {
					defer deleteFiles()
				}()
			})
		}
	}

}

func TestEncryptAndDecryptFailure(t *testing.T) {
	var stdin bytes.Buffer

	for _, elem := range testFailureCase {
		stdin.Write(keytest)
		err := run(elem.cmd, &stdin)
		t.Run(elem.name, func(t *testing.T) {
			if err == nil {
				t.Fatal("error should be occured, but it didn't")
			}

			if err := isFileExist(elem.file[0]); err == nil {
				t.Fatalf("failure: want: %v file is not found\nactual : %v", elem.file[0], err)
			}

			if err := isFileExist(elem.file[1]); err == nil {
				t.Fatalf("failure: want: %v file is not found\nactual : %v", elem.file[1], err)
			}
		})
	}
}