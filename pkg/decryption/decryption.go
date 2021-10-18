package decryption

import (
  "crypto/aes"
	"crypto/cipher"
)

func DecryptAes(cipherData, key []byte) ([]byte, error) {
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