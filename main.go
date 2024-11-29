package main

import (
	file_cryptor "cy-py-go-files/file_cryptor"
	"fmt"
	"sync"
	"time"
)

type CryptorInfo struct {
	// Add fields for your CryptorInfo struct here
}

func fsGetCryptorInfo(wg *sync.WaitGroup, threadID string) (map[string]interface{}, error) {
	defer wg.Done()
	file_test := "\\\\192.168.18.36\\disk1\\data.mp4-version-1"
	file_test = "C:\\\\source\\go-src\\cy-go-file\\data-test\\data.mp4-version-1"
	var info map[string]interface{}
	for i := 0; i < 10; i++ {
		startTime := time.Now()
		fs := file_cryptor.FileCryptor{
			FileName: file_test,
		}
		info, err := fs.GetCryptorInfo() // Replace 'fs' with your actual file system or function
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		fmt.Printf("Thread %s: Time taken: %v\n", threadID, time.Since(startTime))
		fmt.Printf("Thread %s: Info: %v\n", threadID, info)
	}

	return info, nil
}
func main() {

	var wg sync.WaitGroup

	// Adjust numWorkers based on your desired parallelism and resource constraints
	numWorkers := 5

	for i := 0; i < numWorkers; i++ {
		threadID := fmt.Sprintf("thread-%d", i)
		wg.Add(1)
		go fsGetCryptorInfo(&wg, threadID)
	}

	wg.Wait() // Wait for all goroutines to finish before exiting
}
