package main

import (
	"log"
	"os"
	"path"
)

const (
	cfgDirDefault string = ".config/let-us-lock"
	saltFileName  string = "salt.txt"
)

func getHomeDir() string {
	usrDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err.Error())
	}

	return usrDir
}

func findCfgDir(usrDir string) string {
	cfgDir := path.Join(usrDir, cfgDirDefault)

	_, err := os.Stat(cfgDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(cfgDir, 0750)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	return path.Join(cfgDir, saltFileName)
}
