package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Iiqbal2000/let-us-lock/crypt"
	fs "github.com/Iiqbal2000/let-us-lock/filesystem"
	"github.com/Iiqbal2000/let-us-lock/randstr"
	"golang.org/x/crypto/scrypt"
)

// parameters for scrypt algorithm except for aesKeyLen and blockSize variables
// https://en.wikipedia.org/wiki/Scrypt
const (
  costFactor int                      = 32768
  blockSizeFactor int                 = 8
  parallelFactor int                  = 1
  blockSize int                       = 32 // AES-256
  saltPath string                     = "salt.txt"
)

func main() {
  if err := run(os.Args, os.Stdin); err != nil {
    io.WriteString(os.Stderr, err.Error())
    io.WriteString(os.Stderr, "\n")
    os.Exit(2)
  }
}

func setCmd(name string) *flag.FlagSet {
  return flag.NewFlagSet(name, flag.ExitOnError)
}

func run(args []string, stdIn io.Reader) error {
  
  if len(args) < 2 {
    return errors.New("you have to include 'encrypt' or 'decrypt' command")
  }

  var (
    passphrase string
    file *string
    output *string
  )

  inputCmd := strings.ToLower(args[1])
  if inputCmd != "encrypt" && inputCmd != "decrypt" {
    return errors.New("you have to include 'encrypt' or 'decrypt' command")
  }

  flagSet := setCmd(inputCmd)
  file = flagSet.String("f", "", "your file path which you want to encrypt/decrypt")
  output = flagSet.String("o", fmt.Sprintf("%s-result", inputCmd), "your file output name")
  if err := flagSet.Parse(args[2:]); err != nil {
    return err
  }

  // get passphrase input
  io.WriteString(os.Stdout, "Enter your password (minimal 8 characters): ")
  buff := bufio.NewReader(stdIn)
  // ReadString will block until the delimiter is entered
	strBuff, err := buff.ReadString('\n')
	if err != nil {
		return errors.New("password is required")
	}
	// remove the delimeter from the string
	passphrase = strings.TrimSuffix(strBuff, "\n")

  // read file
  fileContent, err := fs.ReadFile(*file)
  if err != nil {
    return errors.New("the file is not found")
  }

  // generate salt
  var salt []byte
  if _, err := os.Stat(saltPath); err != nil {
    salt = randstr.Generate(50)
    fs.WriteFile(salt, saltPath)
  } else {
    salt, err = fs.ReadFile(saltPath)
    if err != nil {
      return errors.New("failure when read salt file")
    }
  }

  // generate key from passphrase using scrypt
	key, err := scrypt.Key([]byte(passphrase), salt, costFactor, blockSizeFactor, parallelFactor, blockSize)
  if err != nil {
    return errors.New("the passphrase is not match")
  }

  var outData []byte
  if inputCmd == "encrypt" {
    outData, err = crypt.Encrypt(fileContent, key)
  } else {
    outData, err = crypt.Decrypt(fileContent, key)
  }

  if err != nil {
    return errors.New(err.Error())
  }
  fs.WriteFile(outData, *output)
  
  return nil
}