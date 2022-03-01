package main

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"golang.org/x/crypto/scrypt"
)

// parameters for scrypt algorithm except for aesKeyLen and
// blockSize variables.
var (
	costFactor      int    = 32768 // number of iteration
	blockSizeFactor int    = 8
	parallelFactor  int    = 1
	blockSize       int    = 32 // one of the size blocks that is used AES-256
	saltFile        string = "salt.txt"
	saltLen         int    = 50 // salt length
)

type key []byte

// derive key from passphrase using scrypt kdf
func (k key) derive() ([]byte, error) {
	// remove the delimeter from the string
	passphrase := bytes.TrimSuffix(k, []byte("\n"))

	// check salt file
	var salt []byte
	if _, err := os.Stat(saltFile); err != nil {
		salt = k.generateSalt(saltLen)
		if err = ioutil.WriteFile(saltFile, salt, 0644); err != nil {
			return nil, err
		}
	} else {
		salt, err = ioutil.ReadFile(saltFile)
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

// generateSalt generates random salt
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
