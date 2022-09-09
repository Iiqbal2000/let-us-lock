package main

import (
	"io"
	"os"
)

func main() {
	application := app{
		args:   os.Args,
		input:  os.Stdin,
		output: os.Stdout,
	}

	err := application.run()
	if err != nil {
		io.WriteString(os.Stdout, "\n")
		io.WriteString(os.Stdout, err.Error())
		io.WriteString(os.Stdout, "\n")
		os.Exit(1)
	}
}
