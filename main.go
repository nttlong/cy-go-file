package main

import (
	file_cryptor "cy-py-go-files/file_cryptor"
	"encoding/json"
	"fmt"
	"os"
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
	file_test := "\\\\192.168.18.36\\disk1\\data.mp4-version-1"
	file_test = "C:\\source\\go-src\\cy-go-file\\data-test\\untitled.png"
	fs := file_cryptor.FileCryptor{
		FileName: file_test,
	}
	info, err := fs.GetCryptorInfo()
	if err != nil {
		fmt.Println(err)
		return
	}

	// conver to json text pretty format
	json_text, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(json_text))
	println(fs.GetCryptorInfo())
	reader, err := fs.OpenRead()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(reader)
	content := make([]byte, reader.FileSize)
	reader.Read(content)
	// write content to file by using os.WriteFile()
	file_write_test := "C:\\source\\go-src\\cy-go-file\\data-test\avatar-logo-decrypt.png"
	os.Create(file_write_test)
	os.WriteFile(file_write_test, content, 0644)

}
