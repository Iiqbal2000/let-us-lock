package main

import "flag"

type action struct {
  cmd *flag.FlagSet
  options tOpts
}
type tOpts map[string]interface{}
type tAct func(args []string) (*action, error)

func SetEncryptAct(args []string) (*action, error) {
  flg := flag.NewFlagSet("encrypt", flag.ExitOnError)
  file := flg.String("f", "file-default", "your file path which you want to encrypt/decrypt")
  output := flg.String("o", "encrypt-result", "your file output name")
  act := action{
    cmd: flg,
    options: tOpts{
      "file": file,
      "output" : output,
    },
  }

  if err := flg.Parse(args[2:]); err != nil {
      return nil, err
  }

  return &act, nil
}

func SetDecryptAct(args []string) (*action, error) {
  flg := flag.NewFlagSet("decrypt", flag.ExitOnError)
  file := flg.String("f", "encrypt-result", "your file path which you want to encrypt/decrypt")
  output := flg.String("o", "decrypt-result", "your file output name")
  act := action{
    cmd: flg,
    options: tOpts{
      "file": file,
      "output" : output,
    },
  }

  if err := flg.Parse(args[2:]); err != nil {
      return nil, err
  }

  return &act, nil
}