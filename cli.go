package main

import (
	"fmt"
	"strings"
)

type CliCommand interface {
	Execute(args []string, kdf keyDerivator) error
	Name() string
}

type CliCommands []CliCommand

func (commands CliCommands) GetCommand(commandName string) (CliCommand, error) {
	for _, command := range commands {
		if strings.EqualFold(command.Name(), commandName) {
			return command, nil
		}
	}

	return nil, fmt.Errorf("unknown command: %s", commandName)
}