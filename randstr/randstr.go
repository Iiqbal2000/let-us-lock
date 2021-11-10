package randstr

import (
	"math/rand"
  "time"
	"sync"
)

func Generate(size int) []byte  {
	var salt []byte
	// ASCII range
	min := 32
	max := 127

	var wg sync.WaitGroup
	
	rand.Seed(time.Now().UnixNano())
	wg.Add(1)

	go func(n int)	{
		for i := 0; i < n; i++ {
			random := rand.Intn(max - min + 1) + min
			salt = append(salt, byte(random))
		}
		wg.Done()
	}(size)
	
	wg.Wait()
	
	return salt
}