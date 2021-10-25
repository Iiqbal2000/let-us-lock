package saltgenerator

import (
	"math/rand"
  "time"

	fs "github.com/Iiqbal2000/let-us-lock/pkg/filesystem"
)

func GenerateSalt(size int) []byte  {
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

func ReadSalt(saltPath string) []byte {
	data, _ := fs.ReadFile(saltPath);
	if data == nil {
    rand.Seed(time.Now().UnixNano())
    salt := GenerateSalt(50)
    fs.WriteFile(salt, saltPath)
  }

	return data
}

