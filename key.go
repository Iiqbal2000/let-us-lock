package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"time"

	"golang.org/x/crypto/scrypt"
)

// parameters for scrypt algorithm except for aesKeyLen and
// blockSize variables.
var (
	costFactor      int = 32768 // number of iteration
	blockSizeFactor int = 8
	parallelFactor  int = 1
	blockSize       int = 32 // one of the size blocks that is used AES-256
	cfgDirDefault   string = ".config/let-us-lock"
	cfgFile         string = "config.txt"
	saltLen         int    = 50 // salt length
)

type key []byte

// getCfgPath constructs a path of config dir.
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
	if _, err := os.Stat(cfgPath); err != nil {
		err := os.MkdirAll(cfgPath, 0750)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return cfgPath
}

// derive derives a key from passphrase using scrypt kdf.
func (k key) derive() ([]byte, error) {
	// remove delimiter from the string.
	passphrase := bytes.TrimSuffix(k, []byte("\n"))

	cfgDir := checkCfgDir(getCfgPath())

	var salt []byte

	if _, err := os.Stat(path.Join(cfgDir, cfgFile)); err != nil {
		salt = k.generateSalt(saltLen)
		if err = ioutil.WriteFile(path.Join(cfgDir, cfgFile), salt, 0644); err != nil {
			return nil, err
		}
	} else {
		salt, err = ioutil.ReadFile(path.Join(cfgDir, cfgFile))
		if err != nil {
			return nil, ErrSaltNotFound
		}
	}

	key, err := scrypt.Key([]byte(passphrase), salt, costFactor, blockSizeFactor, parallelFactor, blockSize)
	if err != nil {
		return nil, ErrPassWrong
	}
	return key, nil
}

// generateSalt generates random salt.
func (k key) generateSalt(size int) []byte {
	var salt []byte
	// ASCII range
	min := 32
	max := 127

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < size; i++ {
		// randomize in ascii range
		random := rand.Intn(max-min+1) + min
		salt = append(salt, byte(random))
	}

	return salt
}
