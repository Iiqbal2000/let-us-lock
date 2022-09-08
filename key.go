package main

import (
	"bytes"
	"crypto/rand"
	"os"
	"path"

	"golang.org/x/crypto/scrypt"
)

// parameters for scrypt algorithm except for aesKeyLen and
// blockSize variables.
var (
	costFactor      int = 32768 // number of iteration
	blockSizeFactor int = 8
	parallelFactor  int = 1
	blockSize       int = 32 // one of the size blocks that is used AES-256
	saltLen         int = 50 // salt length
)

// key contains a passphrase that hashed.
type key struct {
	passphrase []byte
	hashResult []byte
}

// derive derives a key from passphrase using scrypt kdf.
func (k key) derive() ([]byte, error) {
	// remove delimiter from the string.
	passphrase := bytes.TrimSuffix(k, []byte("\n"))

	cfgDirPath := checkCfgDir(getCfgPath())

	var salt []byte

	filePath := path.Join(cfgDirPath, cfgFile)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		salt = k.generateSalt(saltLen)
		if err = os.WriteFile(filePath, salt, 0644); err != nil {
			return nil, err
		}

	} else {
		salt, err = os.ReadFile(filePath)
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
	salt := make([]byte, size)
	rand.Read(salt)

	return salt
}
