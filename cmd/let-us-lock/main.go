package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/Iiqbal2000/let-us-lock/pkg/aes"
	fs "github.com/Iiqbal2000/let-us-lock/pkg/filesystem"
	"github.com/Iiqbal2000/let-us-lock/pkg/randstr"
	"golang.org/x/crypto/scrypt"
)

const (
  N int                 = 32768 // cpu cost
  r int                 = 8 // memory cost
  p int                 = 1 // parallelization
  keyLen int            = 32 // byte key length for AES-256
  saltPath string       = "salt.txt"
)

func main() {
  var mode = flag.String("m", "encrypt", "encrypt/decrypt")
  var file = flag.String("f", "", "your file path which you want to encrypt/decrypt")
  var output = flag.String("o", "", "your file output name")
  var passphrase string

  // parsing input flag
  flag.Parse()

  fmt.Print("Enter your password (minimal 8 characters): ")
  // get passphrase input
  fmt.Scanf("%s", &passphrase)
  
  // generate key from passphrase using scrypt
	key, err := scrypt.Key([]byte(passphrase), randstr.Read(saltPath), N, r, p, keyLen)
  if err != nil {
    log.Fatal(err.Error())
  }

  // read file
  data, err := fs.ReadFile(*file)
  if err != nil {
    log.Fatal(err.Error())
  }
  
  if inputMode := strings.ToLower(*mode); inputMode == "encrypt" {
    cipherData, err := aes.Encrypt(data, key)
    if err != nil {
      log.Fatal(err.Error())
    }
    fs.WriteFile(cipherData, *output)

  } else if inputMode == "decrypt" {
    plainData, err := aes.Decrypt(data, key)
    if err != nil {
      log.Fatal(err.Error())
    }
    fs.WriteFile(plainData, *output)

  } else {
    fmt.Println("Invalid mode")
  }
}