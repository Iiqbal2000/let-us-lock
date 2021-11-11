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
	var mtx sync.Mutex

	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func() {
			mtx.Lock()
			random := rand.Intn(max - min + 1) + min
			salt = append(salt, byte(random))
			mtx.Unlock()
			wg.Done()
		}()
	}
	
	wg.Wait()
	
	return salt
}