package main

import (
	"testing"
	"bytes"
	"os"
)

var key = []byte("passphrasetesti\n")

var testSuccessCase = []struct{
	name string
	cmd []string
	expectFile []string
} {
		{
			name: "encrypt",
			cmd: []string{"main.go", "encrypt", "-f", "kitten.png", "-o", "cipherfile"},
			expectFile: []string{"salt.txt", "cipherfile"},
		},
		{	
			name: "decrypt",
			cmd: []string{"main.go", "decrypt", "-f", "cipherfile", "-o", "result.png"},
			expectFile: []string{"salt.txt", "result.png"},
		},
}

var testFailureCase = []struct{
	name string
	cmd []string
	file []string
} {
	{
		name: "error",
		cmd: []string{"main.go", "ksiwn", "-f", "kitten.png", "-o", "cipherfile"},
		file: []string{"salt.txt", "cipherfile"},
	},
	{
		name: "error",
		cmd: []string{"main.go", "encrypt"},
		file: []string{"salt.txt", "cipherfile"},
	},
	{
		name: "error",
		cmd: []string{"main.go", "-f", "kitten.png", "-o", "cipherfile"},
		file: []string{"salt.txt", "cipherfile"},
	},
	{
		name: "error",
		cmd: []string{"main.go", "encrypt", "-o", "cipherfile"},
		file: []string{"salt.txt", "cipherfile"},
	},
}

func TestEncryptAndDecrypt(t *testing.T) {
	var stdin bytes.Buffer
	
	for _, val := range testSuccessCase {
		stdin.Write(key)
		run(val.cmd, &stdin)

		if val.name == "encrypt" {
			t.Run(val.name, func(t *testing.T) {
				if _, err := os.Stat(val.expectFile[0]); err != nil {
					t.Fatalf("Got error %s file is not found. stack : %s", val.expectFile[0], err)
				}		
				if _, err := os.Stat(val.expectFile[1]); err != nil {
					t.Fatalf("Got error %s file is not found. stack : %s", val.expectFile[1], err)
				}
			})
		} else if val.name == "decrypt" {
			t.Run(val.name, func(t *testing.T) {
				if _, err := os.Stat(val.expectFile[1]); err != nil {
					t.Fatalf("Got error %s file is not found. stack : %s", val.expectFile[1], err)
				}

				func() {
					defer func() {
						if err := os.Remove("salt.txt"); err != nil {
							t.Fatal(err)
						}
						if err := os.Remove("cipherfile"); err != nil {
							t.Fatal(err)
						}
						if err := os.Remove("result.png"); err != nil {
							t.Fatal(err)
						}
					}()
				}()
			})
		}
	}

}

func TestEncryptAndDecryptFailure(t *testing.T) {
	var stdin bytes.Buffer

	for _, elem := range testFailureCase {
		stdin.Write(key)
		err := run(elem.cmd, &stdin)
		t.Run(elem.name, func(t *testing.T) {
			if err == nil {
				t.Fatal("error should appear, but it didn't")
			}

			if _, err := os.Stat(elem.file[0]); err == nil {
				t.Fatalf("Got error %s file is found. stack : %s", elem.file[0], err)
			}

			if _, err := os.Stat(elem.file[1]); err == nil {
				t.Fatalf("Got error %s file is found. stack : %s", elem.file[1], err)
			}
		})
	}
}