package main

import (
	"testing"
	"bytes"
	"os"
)

var key = []byte("passphrasetesti\n")

var runTestCase = []struct{
	name string
	cmd []string
	expectFile []string
	err bool
} {
		{
			name: "encrypt",
			cmd: []string{"main.go", "encrypt", "-f", "../../kitten.png", "-o", "cipherfile"},
			expectFile: []string{"salt.txt", "cipherfile"},
			err: false,
		},
		{	
			name: "decrypt",
			cmd: []string{"main.go", "decrypt", "-f", "cipherfile", "-o", "result.png"},
			expectFile: []string{"salt.txt", "result.png"},
			err: false,
		},
		{
			name: "error",
			cmd: []string{"main.go", "ksiwn", "-f", "../../kitten.png", "-o", "cipherfile"},
			expectFile: []string{"salt.txt", "cipherfile"},
			err: true,
		},
		{
			name: "error",
			cmd: []string{"main.go", "encrypt"},
			expectFile: []string{"salt.txt", "cipherfile"},
			err: true,
		},
		{
			name: "error",
			cmd: []string{"main.go", "-f", "../../kitten.png", "-o", "cipherfile"},
			expectFile: []string{"salt.txt", "cipherfile"},
			err: true,
		},
		{
			name: "error",
			cmd: []string{"main.go", "encrypt", "-o", "cipherfile"},
			expectFile: []string{"salt.txt", "cipherfile"},
			err: true,
		},
}

func TestRun(t *testing.T) {
	var stdin bytes.Buffer
	
	for _, val := range runTestCase {
		stdin.Write(key)
		err := run(val.cmd, &stdin)
		
		if val.name == "encrypt" && !val.err {
			if fileInfo, err := os.Stat(val.expectFile[0]); err != nil {
				t.Fatalf("Got %s, expected: %s", fileInfo.Name(), val.expectFile[0])
			}		
			if fileInfo, err := os.Stat(val.expectFile[1]); err != nil {
				t.Fatalf("Got %s, expected: %s", fileInfo.Name(), val.expectFile[1])
			}
		} else if val.name == "decrypt" && !val.err {
			if fileInfo, err := os.Stat(val.expectFile[1]); err != nil {
				t.Fatalf("Got %s, expected: %s", fileInfo.Name(), val.expectFile[1])
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

		} else {
			if err == nil {
				t.Fatal("an error is not found")
			}

			if fileInfo, err := os.Stat(val.expectFile[0]); err == nil {
				t.Fatalf("Got %s, expected: %s", fileInfo.Name(), "no file")
			}

			if fileInfo, err := os.Stat(val.expectFile[1]); err == nil {
				t.Fatalf("Got %s, expected: %s", fileInfo.Name(), "no file")
			}
		}
	}

}