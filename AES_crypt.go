package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Encrypt encrypt plaintext to ciphertext using AES.
func Encrypt(plainData, key []byte) ([]byte, error) {
  gcm, err := createCipherBlock(key)
  if err != nil { 
    return nil, err
  }

  // generate nonce (number used once).
	// nonce slice length must be from gcm.
  nonce := make([]byte, gcm.NonceSize())
  _, err = io.ReadFull(rand.Reader, nonce)
	if err != nil { 
    return nil, err
  }

  cipherData := gcm.Seal(
		nonce,
		nonce,
		plainData,
		nil,
  )
  
  return cipherData, nil
}

// Decrypt decrypt ciphertext to plaintext using AES.
func Decrypt(cipherData, key []byte) ([]byte, error) {
  gcm, err := createCipherBlock(key)
  if err != nil { 
    return nil, err
  }

  // get nonce
  nonceSize := gcm.NonceSize()
  nonce, cipherText := cipherData[:nonceSize], cipherData[nonceSize:]

  plainData, err := gcm.Open(nil, nonce, cipherText, nil)
  if err != nil { 
    return nil, err
  }
  
  return plainData, nil
}

func createCipherBlock(key []byte) (cipher.AEAD, error) {
  block, err := aes.NewCipher(key)
  if err != nil { 
    return nil, err
  }
  
  return cipher.NewGCM(block)
}