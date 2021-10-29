package randstr

import (
	"math/rand"
  "time"
	"os"
	"sync"

	fs "github.com/Iiqbal2000/let-us-lock/filesystem"
)

func GenerateAndSave(size int, filename string) []byte  {
	var salt []byte
	// ASCII range
	min := 32
	max := 127

	if _, err := os.Stat(filename); err != nil {
		var wg sync.WaitGroup
		rand.Seed(time.Now().UnixNano())
		wg.Add(1)
    func(n int)	{
			for i := 0; i < n; i++ {
				random := rand.Intn(max - min + 1) + min
				salt = append(salt, byte(random))
			}
			wg.Done()
		}(size)
		wg.Wait()
    fs.WriteFile(salt, filename)
	}

	data, _ := fs.ReadFile(filename);
	
	return data
}

// func Read(path string) []byte {
// 	if _, err := os.Stat(path); err != nil {
// 		rand.Seed(time.Now().UnixNano())
//     	salt := Generate(50)
//     	fs.WriteFile(salt, path)
// 	}

// 	data, _ := fs.ReadFile(path);
	
// 	return data
// }

