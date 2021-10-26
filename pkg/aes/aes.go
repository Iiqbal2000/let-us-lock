package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Algoritme AES mengubah blok plaintext 128-bit menjadi blok ciphertext berukuran 128 bit.
func Encrypt(plainData, key []byte) ([]byte, error) {
	// create cipher block
  block, err := aes.NewCipher(key)
  if err != nil { 
    return nil, err
  }
  
  gcm, err := cipher.NewGCM(block)
  if err != nil { 
    return nil, err
  }

  // generate nonce (number used once)
	// nonce slice length must be from gcm
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

func Decrypt(cipherData, key []byte) ([]byte, error) {
  block, err := aes.NewCipher(key)
  if err != nil { 
    return nil, err
  }

  gcm, err := cipher.NewGCM(block)
  if err != nil { 
    return nil, err
  }

  // get nonce
  nonceSize := gcm.NonceSize()
  nonce, cipherText := cipherData[:nonceSize], cipherData[nonceSize:]

  // decrypt
  plainData, err := gcm.Open(nil, nonce, cipherText, nil)
  
  if err != nil { 
    return nil, err
  }
  return plainData, nil
}