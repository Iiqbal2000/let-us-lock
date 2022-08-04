package main

import (
	"errors"
	"io"
	"os"
	"syscall"

	"golang.org/x/term"
)

var (
	ErrCmd          = errors.New("you have to include 'encrypt' or 'decrypt' command")
	ErrPassWrong    = errors.New("the password/passphrase is invalid")
	ErrPassTooShort = errors.New("the password/passphrase too short (MIN 8 characters)")
	ErrPassTooLong  = errors.New("the password/passphrase too long (MAX 64 characters)")
	ErrPassNotFound = errors.New("the password/passphrase is required")
	ErrFileNotFound = errors.New("the file is not found")
	ErrSaltNotFound = errors.New("failure when reading a salt file")
)

func main() {
	if err := run(os.Args, os.Stdin, os.Stdout, true); err != nil {
		io.WriteString(os.Stdout, err.Error())
		io.WriteString(os.Stdout, "\n")
		os.Exit(1)
	}
}

type cryptHandler func(plainData, key []byte) ([]byte, error)

func run(args []string, stdIn, stdOut io.ReadWriter, hidePassword bool) error {
	if len(args) < 2 {
		return ErrCmd
	}

	commands := CliCommands{
		newEncryptCmd(cryptHandler(Encrypt)),
		newDecryptCmd(cryptHandler(Decrypt)),
	}

	// get a command
	cmd, err := commands.GetCommand(args[1])
	if err != nil {
		return err
	}

	err = cmd.Validate(args)
	if err != nil {
		return err
	}

	// get a passphrase
	io.WriteString(stdOut, "Enter your password (min 8 characters and max 64 characters): ")

	var rawPassphrase []byte

	if hidePassword {
		rawPassphrase, err = term.ReadPassword(int(syscall.Stdin))
		io.WriteString(stdOut, "\n")
	} else {
		rawPassphrase, err = io.ReadAll(stdIn)
	}

	if err != nil {
		return ErrPassNotFound
	}

	// Validate the passphrase
	if len(rawPassphrase) < 8 {
		return ErrPassTooShort
	} else if len(rawPassphrase) > 64 {
		return ErrPassTooLong
	}

	err = cmd.Execute(key(rawPassphrase))
	if err != nil {
		return err
	}

	return nil
}
