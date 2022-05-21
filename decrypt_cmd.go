package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type decryptCmd struct {
	flagSet    *flag.FlagSet
	inputPath  string
	outputPath string
	decrypt    cryptHandler
}

func newDecryptCmd(decryptH cryptHandler) *decryptCmd {
	decryptcmd := &decryptCmd{
		flagSet: flag.NewFlagSet("decrypt", flag.ExitOnError),
		decrypt: decryptH,
	}

	decryptcmd.flagSet.StringVar(&decryptcmd.inputPath, "f", "encrypt-result", "your file path which you want to decrypt")
	decryptcmd.flagSet.StringVar(&decryptcmd.outputPath, "o", "decrypt-result", "your file output name")
	decryptcmd.flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "USAGE:")
		fmt.Fprintln(os.Stderr, "   decrypt -f [your file] -o [your new file]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "COMMANDS:")
		fmt.Fprintln(os.Stderr, "   decrypt - to decrypt a file")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "OPTIONS:")
		decryptcmd.flagSet.PrintDefaults()
	}
	return decryptcmd
}

func (c *decryptCmd) Execute(args []string, kdf key) error {
	if err := c.flagSet.Parse(args[2:]); err != nil {
		return err
	}

	// read and check file
	data, err := os.ReadFile(c.inputPath)
	if err != nil {
		return ErrFileNotFound
	}

	key, err := kdf.derive()
	if err != nil {
		return err
	}

	var plaintext []byte
	plaintext, err = c.decrypt(data, key)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if err = os.WriteFile(c.outputPath, plaintext, 0644); err != nil {
		return err
	}
	return nil
}

func (c decryptCmd) Name() string {
	return c.flagSet.Name()
}
