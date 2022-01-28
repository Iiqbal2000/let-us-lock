package main

import (
	"math/rand"
	"time"
)

func GenerateSalt(size int) []byte  {
	var salt []byte
	// ASCII range
	min := 32
	max := 127

	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < size; i++ {
		// randomize in ascii range
		random := rand.Intn(max - min + 1) + min
		salt = append(salt, byte(random))
	}
		
	return salt
}