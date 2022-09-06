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
	err := run(config{
		args: os.Args,
		stdIn: os.Stdin,
		stdOut: os.Stdout,
	})
	if err != nil {
		io.WriteString(os.Stdout, err.Error())
		io.WriteString(os.Stdout, "\n")
		os.Exit(1)
	}
}

type cryptHandler func(plainData, key []byte) ([]byte, error)

type config struct {
	args []string
	stdIn io.ReadWriter
	stdOut io.ReadWriter
}

func run(conf config) error {
	if len(conf.args) < 2 {
		return ErrCmd
	}

	commands := CliCommands{
		newEncryptCmd(cryptHandler(Encrypt)),
		newDecryptCmd(cryptHandler(Decrypt)),
	}

	// get a command
	cmd, err := commands.GetCommand(conf.args[1])
	if err != nil {
		return err
	}

	err = cmd.Validate(conf.args)
	if err != nil {
		return err
	}

	// get a passphrase
	io.WriteString(conf.stdOut, "Enter your password (min 8 characters and max 64 characters): ")

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

func runForTesting(conf config) error {
	if len(conf.args) < 2 {
		return ErrCmd
	}

	commands := CliCommands{
		newEncryptCmd(cryptHandler(Encrypt)),
		newDecryptCmd(cryptHandler(Decrypt)),
	}

	// get a command
	cmd, err := commands.GetCommand(conf.args[1])
	if err != nil {
		return err
	}

	err = cmd.Validate(conf.args)
	if err != nil {
		return err
	}

	// get a passphrase
	io.WriteString(conf.stdOut, "Enter your password (min 8 characters and max 64 characters): ")

	passphrase, err := catchPassphrase(io.ReadAll(conf.stdIn))
	if err != nil {
		return err
	}

	err = cmd.Execute(key(passphrase))
	if err != nil {
		return err
	}

	return nil
}