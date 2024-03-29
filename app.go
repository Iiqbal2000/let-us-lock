package main

import (
	"io"
	"syscall"

	"golang.org/x/term"
)

const (
	minArgsLen             = 2
	msgForAskingPassphrase = "Enter your password (min 8 characters and max 64 characters): "
)

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

	commands := Commands{
		newEncryptCmd(cryptHandler(Encrypt)),
		newDecryptCmd(cryptHandler(Decrypt)),
	}

	cmdInput := ap.args[1]
	cmd, err := commands.Get(cmdInput)
	if err != nil {
		return err
	}

	err = cmd.Validate(ap.args)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ap.output, msgForAskingPassphrase)
	if err != nil {
		return err
	}

	passphraseInput, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return ErrPassNotFound
	}

	key, err := createKey(passphraseInput)
	if err != nil {
		return err
	}

	err = cmd.Execute(key.HashResult())
	if err != nil {
		return err
	}

	return nil
}

func (ap app) runForTesting() error {
	if len(ap.args) < minArgsLen {
		return ErrCmd
	}

	commands := Commands{
		newEncryptCmd(cryptHandler(Encrypt)),
		newDecryptCmd(cryptHandler(Decrypt)),
	}

	cmdInput := ap.args[1]
	cmd, err := commands.Get(cmdInput)
	if err != nil {
		return err
	}

	err = cmd.Validate(ap.args)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ap.output, msgForAskingPassphrase)
	if err != nil {
		return err
	}

	passphraseInput, err := io.ReadAll(ap.input)
	if err != nil {
		return ErrPassNotFound
	}

	passphrase, err := createKey(passphraseInput)
	if err != nil {
		return err
	}

	err = cmd.Execute(passphrase.HashResult())
	if err != nil {
		return err
	}

	return nil
}
