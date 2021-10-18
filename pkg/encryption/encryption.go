package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Algoritme AES mengubah blok plaintext 128-bit menjadi blok ciphertext berukuran 128 bit.
func EncryptAes(plainData, key []byte) ([]byte, error) {
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