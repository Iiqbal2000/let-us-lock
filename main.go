package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/Iiqbal2000/let-us-lock/crypt"
	fs "github.com/Iiqbal2000/let-us-lock/filesystem"
	"github.com/Iiqbal2000/let-us-lock/randstr"
	"golang.org/x/crypto/scrypt"
)

// parameters for scrypt algorithm except for aesKeyLen and blockSize variables
// https://en.wikipedia.org/wiki/Scrypt
var (
  costFactor int                      = 32768
  blockSizeFactor int                 = 8
  parallelFactor int                  = 1
  blockSize int                       = 32 // AES-256
  saltFile string                     = "salt.txt"
  saltLen int                         = 50
)

var (
  ErrCmd = errors.New("you have to include 'encrypt' or 'decrypt' command")
  ErrPassWrong = errors.New("the passphrase is not match")
  ErrPassNotFound = errors.New("password is required")
  ErrFileNotFound = errors.New("the file is not found")
  ErrSaltNotFound = errors.New("failure when read salt file")
)

func main() {
  if err := run(os.Args, os.Stdin); err != nil {
    io.WriteString(os.Stdout, err.Error())
    io.WriteString(os.Stdout, "\n")
    os.Exit(1)
  }
}

type keyDerivator []byte

func (k keyDerivator) derive() ([]byte, error) {
	// remove the delimeter from the string
	passphrase := bytes.TrimSuffix(k, []byte("\n"))
  
  // check salt
  var salt []byte
  if _, err := os.Stat(saltFile); err != nil {
    salt = randstr.Generate(saltLen)
    fs.WriteFile(salt, saltFile)
  } else {
    salt, err = fs.ReadFile(saltFile)
    if err != nil {
      return nil, ErrSaltNotFound
    }
  }

  // derive key from passphrase using scrypt kdf
	key, err := scrypt.Key([]byte(passphrase), salt, costFactor, blockSizeFactor, parallelFactor, blockSize)
  if err != nil {
    return nil, ErrPassWrong
  }
	return key, nil
}

type cryptHandler func(plainData, key []byte) ([]byte, error)

func run(args []string, stdIn io.Reader) error {
  
  if len(args) < 2 {
    return ErrCmd
  }
  
  commands := CliCommands{newEncryptCmd(cryptHandler(crypt.Encrypt)), newDecryptCmd(cryptHandler(crypt.Decrypt))}
  cmd, err := commands.GetCommand(args[1])
  if err != nil {
    return err
  }
  
  // get passphrase input
  io.WriteString(os.Stdout, "Enter your password (minimal 8 characters): ")
  buff := bufio.NewReader(stdIn)
  // ReadString will block until the delimiter is entered
	rawPassphrase, err := buff.ReadBytes('\n')
	if err != nil {
		return ErrPassNotFound
	}
  
  err = cmd.Execute(args, keyDerivator(rawPassphrase))
  if err != nil {
    return err
  }
  
  return nil
}