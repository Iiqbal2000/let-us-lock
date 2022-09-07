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

const minArgsLen = 2

func main() {
	application := app{
		args:   os.Args,
		input:  os.Stdin,
		output: os.Stdout,
	}

	err := application.run()
	if err != nil {
		io.WriteString(os.Stdout, err.Error())
		io.WriteString(os.Stdout, "\n")
		os.Exit(1)
	}
}

type cryptHandler func(plainData, key []byte) ([]byte, error)

type app struct {
	args   []string
	input  io.ReadWriter
	output io.ReadWriter
}

func (ap app) run() error {
	if len(ap.args) < minArgsLen {
		return ErrCmd
	}

	commands := CliCommands{
		newEncryptCmd(cryptHandler(Encrypt)),
		newDecryptCmd(cryptHandler(Decrypt)),
	}

	cmd, err := commands.GetCommand(ap.args[1])
	if err != nil {
		return err
	}

	err = cmd.Validate(ap.args)
	if err != nil {
		return err
	}

	io.WriteString(ap.output, "Enter your password (min 8 characters and max 64 characters): ")

	passphrase, err := catchPassphrase(term.ReadPassword(int(syscall.Stdin)))
	if err != nil {
		return err
	}

	err = cmd.Execute(key(passphrase))
	if err != nil {
		return err
	}

	return nil
}

func (ap app) runForTesting() error {
	if len(ap.args) < minArgsLen {
		return ErrCmd
	}

	commands := CliCommands{
		newEncryptCmd(cryptHandler(Encrypt)),
		newDecryptCmd(cryptHandler(Decrypt)),
	}

	cmd, err := commands.GetCommand(ap.args[1])
	if err != nil {
		return err
	}

	err = cmd.Validate(ap.args)
	if err != nil {
		return err
	}

	io.WriteString(ap.output, "Enter your password (min 8 characters and max 64 characters): ")

	passphrase, err := catchPassphrase(io.ReadAll(ap.input))
	if err != nil {
		return err
	}

	err = cmd.Execute(key(passphrase))
	if err != nil {
		return err
	}

	return nil
}
