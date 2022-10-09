package main

import (
	"fmt"
	"strings"
)

type Command interface {
	Execute(key []byte) error
	Name() string
	Validate(args []string) error
}

type Commands []Command

func (cmds Commands) Get(cmdIn string) (Command, error) {
	for _, command := range cmds {
		if strings.EqualFold(command.Name(), cmdIn) {
			return command, nil
		}
	}

	return nil, fmt.Errorf("unknown command: %s", cmdIn)
}
