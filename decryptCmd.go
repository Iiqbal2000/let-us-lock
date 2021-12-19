package main

import (
	"flag"
	"fmt"
	"os"

	fs "github.com/Iiqbal2000/let-us-lock/filesystem"
)

type decryptCmd struct {
	flagSet flag.FlagSet
	file string
	output string
	handler cryptHandler
}

func newDecryptCmd(handler cryptHandler) *decryptCmd {
	decryptcmd := &decryptCmd{*flag.NewFlagSet("decrypt", flag.ExitOnError), "", "", handler}
	decryptcmd.flagSet.StringVar(&decryptcmd.file,  "f", "encrypt-result", "your file path which you want to decrypt")
	decryptcmd.flagSet.StringVar(&decryptcmd.output,  "o", "decrypt-result", "your file output name")
	decryptcmd.flagSet.Usage = func () {
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

func (c *decryptCmd) Execute(args []string, kdf keyDerivator) error {
	if err := c.flagSet.Parse(args[2:]); err != nil {
    return err
  }

	// read and check file
	fileContent, err := fs.ReadFile(c.file)
	if err != nil {
		return ErrFileNotFound
	}

	key, err := kdf.derive()
	if err != nil {
    return err
  }

	var outData []byte
	outData, err = c.handler(fileContent, key)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	fs.WriteFile(outData, c.output)
	return nil
}

func (c *decryptCmd) Name() string {
	return c.flagSet.Name()
}