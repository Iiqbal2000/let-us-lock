package main

import (
	"log"
	"os"
	"path"
)

const (
	cfgDirDefault   string = ".config/let-us-lock"
	cfgFile         string = "config.txt"
)

// getCfgPath constructs a path of the config dir.
func getCfgPath() string {
	usr, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err.Error())
	}

	return path.Join(usr, cfgDirDefault)
}

// checkCfgDir checks config dir if the config dir does not
// exist, it will create.
func checkCfgDir(cfgPath string) string {
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		err := os.MkdirAll(cfgPath, 0750)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return cfgPath
}
