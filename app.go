package main

import (
	"io"
	"syscall"

	"golang.org/x/term"
)

const minArgsLen = 2

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

	key, err := catchPassphrase(term.ReadPassword(int(syscall.Stdin)))
	if err != nil {
		return err
	}

	err = key.derive()
	if err != nil {
		return err
	}

	err = cmd.Execute(key.hash())
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

	err = passphrase.derive()
	if err != nil {
		return err
	}

	err = cmd.Execute(passphrase.hash())
	if err != nil {
		return err
	}

	return nil
}
