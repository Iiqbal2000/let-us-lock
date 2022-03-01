package main

import (
	"flag"
	"fmt"
	"os"
	"io/ioutil"
)

type encryptCmd struct {
	flagSet flag.FlagSet
	file string
	output string
	encrypt cryptHandler
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
  data, err := ioutil.ReadFile(c.file)
  if err != nil {
    return ErrFileNotFound
  }

	key, err := kdf.derive()
	if err != nil {
    return err
  }

	var chipertext []byte
	chipertext, err = c.encrypt(data, key)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	
	if err = ioutil.WriteFile(c.output, chipertext, 0644); err != nil { 
    return err
  }
	return nil
}

func (c *encryptCmd) Name() string {
	return c.flagSet.Name()
}