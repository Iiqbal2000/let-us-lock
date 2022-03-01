package main

import (
	"bufio"
	"errors"
	"io"
	"os"
)

var (
	ErrCmd          = errors.New("you have to include 'encrypt' or 'decrypt' command")
	ErrPassWrong    = errors.New("the passphrase is not match")
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

type cryptHandler func(plainData, key []byte) ([]byte, error)

func run(args []string, stdIn io.Reader) error {
	if len(args) < 2 {
		return ErrCmd
	}

	commands := CliCommands{
		newEncryptCmd(cryptHandler(Encrypt)),
		newDecryptCmd(cryptHandler(Decrypt)),
	}

  // get command that put by the user
	cmd, err := commands.GetCommand(args[1])
	if err != nil {
		return err
	}

	// get passphrase from user input
	io.WriteString(os.Stdout, "Enter your password (minimal 8 characters): ")

	buff := bufio.NewReader(stdIn)

	// read passphrase until \n
	rawPassphrase, err := buff.ReadBytes('\n')
	if err != nil {
		return ErrPassNotFound
	}

	err = cmd.Execute(args, key(rawPassphrase))
	if err != nil {
		return err
	}

	return nil
}
