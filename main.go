package main

import (
	file_cryptor "cy-py-go-files/file_cryptor"
	"fmt"
	"time"
)

func main() {
	file_test := "\\\\192.168.18.36\\disk1\\data.mp4-version-1"
	fmt.Println("Hello, world!")
	fs := file_cryptor.FileCryptor{
		FileName: file_test,
	}
	for i := 0; i < 10; i++ {
		// get the time call the function
		star_at := time.Now()
		info, err := fs.GetCryptorInfo()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(time.Since(star_at))
		fmt.Println(info)
	}
}
