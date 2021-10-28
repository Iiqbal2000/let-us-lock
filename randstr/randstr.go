package randstr

import (
	"math/rand"
  "time"
	"os"

	fs "github.com/Iiqbal2000/let-us-lock/filesystem"
)

func Generate(size int) []byte  {
	var salt []byte
  // ASCII range
	min := 32
	max := 127
	
	for i := 0; i < size; i++ {
		random := rand.Intn(max-min+1) + min
		salt = append(salt, byte(random))
	}

	return salt
}

func Read(path string) []byte {
	if _, err := os.Stat(path); err != nil {
		rand.Seed(time.Now().UnixNano())
    	salt := Generate(50)
    	fs.WriteFile(salt, path)
	}

	data, _ := fs.ReadFile(path);
	
	return data
}

