package main

import (
	"flag"
	"fmt"
	"os"
)

type encryptCmd struct {
	flagSet    *flag.FlagSet
	inputPath  string
	outputPath string
	encrypt    cryptHandler
}

func newEncryptCmd(encryptH cryptHandler) *encryptCmd {
	encryptcmd := &encryptCmd{
		flagSet: flag.NewFlagSet("encrypt", flag.ExitOnError),
		encrypt: encryptH,
	}

	encryptcmd.flagSet.StringVar(&encryptcmd.inputPath, "f", "file", "your file path which you want to encrypt")
	encryptcmd.flagSet.StringVar(&encryptcmd.outputPath, "o", "encrypt-result", "your file output name")
	
	encryptcmd.flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "USAGE:")
		fmt.Fprintln(os.Stderr, "   encrypt -f [your file] -o [your new file]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "COMMANDS:")
		fmt.Fprintln(os.Stderr, "   encrypt - to encrypt a file")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "OPTIONS:")
		encryptcmd.flagSet.PrintDefaults()
	}

	return encryptcmd
}

func (c *encryptCmd) Validate(args []string) error {
	if err := c.flagSet.Parse(args[2:]); err != nil {
		return err
	}

	_, err := os.Stat(c.inputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotFound
		}
		return err
	}

	return nil
}

func (c *encryptCmd) Execute(key []byte) error {
	data, err := os.ReadFile(c.inputPath)
	if err != nil {
		return ErrFileNotFound
	}

	chipertext, err := c.encrypt(data, key)
	if err != nil {
		return err
	}

	if err = os.WriteFile(c.outputPath, chipertext, 0644); err != nil {
		return err
	}
	return nil
}

func (c encryptCmd) Name() string {
	return c.flagSet.Name()
}
