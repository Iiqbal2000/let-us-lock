package main

import (
	"bufio"
	"errors"
	"flag"
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
var (
  costFactor int                      = 32768
  blockSizeFactor int                 = 8
  parallelFactor int                  = 1
  blockSize int                       = 32 // AES-256
  saltPath string                     = "salt.txt"
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
    io.WriteString(os.Stderr, err.Error())
    io.WriteString(os.Stderr, "\n")
    os.Exit(2)
  }
}

type action struct {
  cmd *flag.FlagSet
  options tOpts
}
type tOpts map[string]interface{}
type tAct func(args []string) (*action, error)

func SetEncryptAct(args []string) (*action, error) {
  flg := flag.NewFlagSet("encrypt", flag.ExitOnError)
  file := flg.String("f", "", "your file path which you want to encrypt/decrypt")
  output := flg.String("o", "encrypt-result", "your file output name")
  act := action{
    cmd: flg,
    options: tOpts{
      "file": file,
      "output" : output,
    },
  }

  if err := flg.Parse(args[2:]); err != nil {
      return nil, err
  }

  return &act, nil
}

func SetDecryptAct(args []string) (*action, error) {
  flg := flag.NewFlagSet("decrypt", flag.ExitOnError)
  file := flg.String("f", "", "your file path which you want to encrypt/decrypt")
  output := flg.String("o", "decrypt-result", "your file output name")
  act := action{
    cmd: flg,
    options: tOpts{
      "file": file,
      "output" : output,
    },
  }

  if err := flg.Parse(args[2:]); err != nil {
      return nil, err
  }

  return &act, nil
}

var actions = map[string]tAct{
  "encrypt": SetEncryptAct,
  "decrypt": SetDecryptAct,
}

func run(args []string, stdIn io.Reader) error {
  
  if len(args) < 2 {
    return ErrCmd
  }

  var (
    passphrase string
    file *string
    output *string
  )

  inputCmd := strings.ToLower(args[1])
  
  if _, ok := actions[inputCmd]; !ok {
    return ErrCmd
  }

  cmd, err := actions[inputCmd](args)
  if err != nil {
    return err
  }
  
  file = cmd.options["file"].(*string)
  output = cmd.options["output"].(*string)

  // get passphrase input
  io.WriteString(os.Stdout, "Enter your password (minimal 8 characters): ")
  buff := bufio.NewReader(stdIn)
  // ReadString will block until the delimiter is entered
	strBuff, err := buff.ReadString('\n')
	if err != nil {
		return ErrPassNotFound
	}
	// remove the delimeter from the string
	passphrase = strings.TrimSuffix(strBuff, "\n")

  // read file
  fileContent, err := fs.ReadFile(*file)
  if err != nil {
    return ErrFileNotFound
  }

  // generate salt
  var salt []byte
  if _, err := os.Stat(saltPath); err != nil {
    salt = randstr.Generate(saltLen)
    fs.WriteFile(salt, saltPath)
  } else {
    salt, err = fs.ReadFile(saltPath)
    if err != nil {
      return ErrSaltNotFound
    }
  }

  // generate key from passphrase using scrypt
	key, err := scrypt.Key([]byte(passphrase), salt, costFactor, blockSizeFactor, parallelFactor, blockSize)
  if err != nil {
    return ErrPassWrong
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