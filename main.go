package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"syscall"

	"golang.org/x/term"
)

var (
	ErrCmd          = errors.New("you have to include 'encrypt' or 'decrypt' command")
	ErrPassWrong    = errors.New("the passphrase is not match")
	ErrPassNotFound = errors.New("password is required")
	ErrFileNotFound = errors.New("the file is not found")
	ErrSaltNotFound = errors.New("failure when read salt file")
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

	// get command that put by the user
	cmd, err := commands.GetCommand(args[1])
	if err != nil {
		return err
	}

	// get passphrase from user input
	io.WriteString(stdOut, "Enter your password (min 8 characters and max 64 characters): ")
	
	var rawPassphrase []byte

	if hidePassword {
		rawPassphrase, err = term.ReadPassword(int(syscall.Stdin))
		io.WriteString(stdOut, "\n")
	} else {
		buff := bufio.NewReader(stdIn)
		rawPassphrase, err = buff.ReadBytes('\n')
	}

	if err != nil {
		return ErrPassNotFound
	}

	err = cmd.Execute(args, key(rawPassphrase))
	if err != nil {
		return err
	}

	return nil
}
