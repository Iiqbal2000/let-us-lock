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

type key struct {
	passphrase []byte
	hashResult []byte
}

const (
	minPassphraseLength = 8
	maxPassphraseLength = 64
)

func createKey(passphraseIn []byte, err error) (*key, error) {
	if err != nil {
		return &key{}, ErrPassNotFound
	}

	passphrase := &key{
		passphrase: passphraseIn,
	}

	passphrase, err = passphrase.clean().validate()
	if err != nil {
		return &key{}, err
	}

	return passphrase, nil
}

func (k key) validate() (*key, error) {
	if len(k.passphrase) < minPassphraseLength {
		return &key{}, ErrPassTooShort
	} else if len(k.passphrase) > maxPassphraseLength {
		return &key{}, ErrPassTooLong
	}

	return &k, nil
}

func (k key) clean() *key {
	return &key{
		passphrase: bytes.TrimSuffix(k.passphrase, []byte("\n")),
	}
}

func (k *key) derive() error {
	salt, err := k.getSalt()
	if err != nil {
		return err
	}

	hashResult, err := scrypt.Key(
		k.passphrase,
		salt,
		costFactor,
		blockSizeFactor,
		parallelFactor,
		blockSize,
	)
	if err != nil {
		return ErrPassWrong
	}

	k.hashResult = hashResult

	return nil
}

func (k key) getSalt() ([]byte, error) {
	cfgDirPath := hasCfgDir(getCfgPath())

	var salt []byte

	filePath := path.Join(cfgDirPath, cfgFile)

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
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

	return salt, nil
}

// func (k key) readExistingSalt() {

// }

// generateSalt generates random salt.
func (k key) generateSalt(size int) []byte {
	salt := make([]byte, size)
	rand.Read(salt)

	return salt
}

func (k key) HashResult() []byte {
	return k.hashResult
}
