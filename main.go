package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
  "flag"
  "fmt"

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
    io.WriteString(os.Stdout, err.Error())
    io.WriteString(os.Stdout, "\n")
    os.Exit(1)
  }
}

func run(args []string, stdIn io.Reader) error {
  
  if len(args) < 2 {
    return ErrCmd
  }

  var (
    file string
    output string
    cmdLine *flag.FlagSet
  )

  if args[1] == "encrypt" {
    cmdLine = flag.NewFlagSet("encrypt", flag.ExitOnError)
  } else if args[1] == "decrypt" {
    cmdLine = flag.NewFlagSet("decrypt", flag.ExitOnError)
  } else {
    return ErrCmd
  }

  cmdLine.StringVar(&file, "f", "encrypt-result", "your file path which you want to encrypt or decrypt")
  cmdLine.StringVar(&output, "o", "decrypt-result", "your file output name")
  
  cmdLine.Usage = func () {
    fmt.Fprintln(os.Stderr, "USAGE:")
    fmt.Fprintf(os.Stderr, "   %s -f [your file] -o [your new file]\n", args[1])
    fmt.Fprintln(os.Stderr, "")
    fmt.Fprintln(os.Stderr, "COMMANDS:")
    fmt.Fprintf(os.Stderr, "   %s - to %s a file\n", args[1], args[1])
    fmt.Fprintln(os.Stderr, "")
    fmt.Fprintln(os.Stderr, "OPTIONS:")
    cmdLine.PrintDefaults()
  }

  if err := cmdLine.Parse(args[2:]); err != nil {
    return err
  }

  // get passphrase input
  io.WriteString(os.Stdout, "Enter your password (minimal 8 characters): ")
  buff := bufio.NewReader(stdIn)
  // ReadString will block until the delimiter is entered
	contentBuff, err := buff.ReadString('\n')
	if err != nil {
		return ErrPassNotFound
	}
	// remove the delimeter from the string
	passphrase := strings.TrimSuffix(contentBuff, "\n")

  // read file
  fileContent, err := fs.ReadFile(file)
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
  inputCmd := strings.ToLower(args[1])

  if inputCmd == "encrypt" {
    outData, err = crypt.Encrypt(fileContent, key)
    if err != nil {
      return errors.New(err.Error())
    }
  } else {
    outData, err = crypt.Decrypt(fileContent, key)
    if err != nil {
      return errors.New(err.Error())
    }
  }

  fs.WriteFile(outData, output)
  
  return nil
}