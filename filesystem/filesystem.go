package filesystem

import (
	"os"
	"sync"
)

type chunk map[string]int
type buffchunk []byte

const buffsize = 100

func ReadFile(path string) ([]byte, error) {
  file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

  defer file.Close()

  fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

  filesize := int(fileinfo.Size())
	
	worker := filesize / buffsize
	chunks := make([]chunk, worker)

  for i := 0; i < worker; i++ {
		chunks[i] = chunk{"bufsize": buffsize, "offset": buffsize * i}
	}

  if remainder := filesize % buffsize; remainder != 0 {
		workerremainder := chunk{"bufsize": remainder, "offset": worker * buffsize}
		worker++
		chunks = append(chunks, workerremainder)
	}
  
  var wg sync.WaitGroup
  var mx sync.Mutex

  buffchunks := make([]buffchunk, worker)

  wg.Add(worker)

  errAlert := false
  var errstack error

  for i := 0; i < worker; i++ {
		go func(chk chunk, x int) {
			defer wg.Done()
			mx.Lock()
			buff := make([]byte, chk["bufsize"])		
			_, err := file.ReadAt(buff, int64(chk["offset"]))
			buffchunks[x] = buff
			chk = nil
			if err != nil {
        errstack = err
        errAlert = true
				return
			}

			mx.Unlock()
		}(chunks[i], i)
	}

	wg.Wait()

  if errAlert {
    return nil, errstack
  }

	var result []byte

	for _, elem := range buffchunks {
		result = append(result, elem...)
	}

  return result, nil
}

func WriteFile(data []byte, fileName string) (string, error) {
  f, err := os.Create(fileName)

  if err != nil {
      return "", nil
  }

  defer f.Close()
	f.Write(data)

  return fileName, nil
}