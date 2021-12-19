package main

import (
	"flag"
	"fmt"
	"os"

	fs "github.com/Iiqbal2000/let-us-lock/filesystem"
)

type encryptCmd struct {
	flagSet flag.FlagSet
	file string
	output string
	handler cryptHandler
}

func newEncryptCmd(handler cryptHandler) *encryptCmd {
	encryptcmd := &encryptCmd{*flag.NewFlagSet("encrypt", flag.ExitOnError), "", "", handler}
	encryptcmd.flagSet.StringVar(&encryptcmd.file,  "f", "file", "your file path which you want to encrypt")
	encryptcmd.flagSet.StringVar(&encryptcmd.output,  "o", "encrypt-result", "your file output name")
	encryptcmd.flagSet.Usage = func () {
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

func (c *encryptCmd) Execute(args []string, kdf keyDerivator) error {
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
	// todo encrypt handler
	outData, err = c.handler(fileContent, key)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	fs.WriteFile(outData, c.output)
	return nil
}

func (c *encryptCmd) Name() string {
	return c.flagSet.Name()
}