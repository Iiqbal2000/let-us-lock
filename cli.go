package main

import (
	"fmt"
	"strings"
)

type CliCommand interface {
	Execute(key []byte) error
	Name() string
	Validate(args []string) error
}

type CliCommands []CliCommand

func (commands CliCommands) GetCommand(cmdIn string) (CliCommand, error) {
	for _, command := range commands {
		if strings.EqualFold(command.Name(), cmdIn) {
			return command, nil
		}
	}

	return nil, fmt.Errorf("unknown command: %s", cmdIn)
}
