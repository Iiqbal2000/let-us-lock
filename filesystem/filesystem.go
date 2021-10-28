package filesystem

import (
  "io/ioutil"
)

func ReadFile(path string) ([]byte, error) {
  data, err := ioutil.ReadFile(path)
  if err != nil { 
    return nil, err
  }

  return data, nil
}

func WriteFile(data []byte, fileName string) (string, error) {
  err := ioutil.WriteFile(fileName, data, 0644)
  if err != nil { 
    return "", err
  }
  return fileName, nil
}