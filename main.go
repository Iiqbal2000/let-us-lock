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

const (
  N int                 = 32768 // cpu cost
  r int                 = 8 // memory cost
  p int                 = 1 // parallelization
  keyLen int            = 32 // byte key length for AES-256
  saltPath string       = "salt.txt"
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

  switch inputCmd {
    case "encrypt", "decrypt":
      flagSet := setCmd(inputCmd)
      file = flagSet.String("f", "", "your file path which you want to encrypt/decrypt")
      output = flagSet.String("o", fmt.Sprintf("%s-result", inputCmd), "your file output name")
      if err := flagSet.Parse(args[2:]); err != nil {
        return err
      }
    default:
      return errors.New("you have to include 'encrypt' or 'decrypt' command")
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
  
  // generate key from passphrase using scrypt
	key, err := scrypt.Key([]byte(passphrase), randstr.GenerateAndSave(50, saltPath), N, r, p, keyLen)
  if err != nil {
    return errors.New("the passphrase is not match")
  }

  switch inputCmd {
    case "encrypt":
      cipherData, err := crypt.Encrypt(fileContent, key)
      if err != nil {
        return errors.New(err.Error())
      }
      fs.WriteFile(cipherData, *output)
    case "decrypt":
      plainData, err := crypt.Decrypt(fileContent, key)
      if err != nil {
        return errors.New(err.Error())
      }
      fs.WriteFile(plainData, *output)
  }
  
  return nil
}