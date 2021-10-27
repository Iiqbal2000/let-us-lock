package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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
  if err := run(os.Args, os.Stdin); err != nil {
		log.Fatalln(err)
	}
}

func run(args []string, stdIn io.Reader) error {
  
  if len(args) < 2 {
    return errors.New("expected 'encrypt' or 'decrypt' subcommands")
  }

  var (
    passphrase string
    file *string
    output *string
  )

  switch strings.ToLower(args[1]) {
    case "encrypt":
      encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
      file = encryptCmd.String("f", "", "your file path which you want to encrypt/decrypt")
      output = encryptCmd.String("o", "cipherfile", "your file output name")
      if err := encryptCmd.Parse(args[2:]); err != nil {
        return err
      }
    case "decrypt":
      decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
      file = decryptCmd.String("f", "", "your file path which you want to encrypt/decrypt")
      output = decryptCmd.String("o", "result", "your file output name")
      if err := decryptCmd.Parse(args[2:]); err != nil {
        return err
      }
    default:
      return errors.New("expected 'encrypt' or 'decrypt' subcommands")
  }

  // get passphrase input
  fmt.Print("Enter your password (minimal 8 characters): ")
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
  
  // generate key from passphrase using scrypt
	key, err := scrypt.Key([]byte(passphrase), randstr.Read(saltPath), N, r, p, keyLen)
  if err != nil {
    return errors.New("the passphrase is not match")
  }

  switch strings.ToLower(args[1]) {
    case "encrypt":
      cipherData, err := aes.Encrypt(fileContent, key)
      if err != nil {
        return errors.New(err.Error())
      }
      fs.WriteFile(cipherData, *output)
    case "decrypt":
      plainData, err := aes.Decrypt(fileContent, key)
      if err != nil {
        return errors.New(err.Error())
      }
      fs.WriteFile(plainData, *output)
    default:
      return errors.New("command are required")
  }
  
  return nil
}