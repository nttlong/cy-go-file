package file_cryptor

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
)

type FileCryptor struct {
	FileName string
	FileSize int64 // Use int64 for large file sizes
}

func invertBufferParallel(buffer []byte) {
	numWorkers := runtime.NumCPU()                           // Get the number of CPU cores
	chunkSize := (len(buffer) + numWorkers - 1) / numWorkers // Calculate chunk size, rounding up

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			// Calculate start and end indices for this worker's chunk
			start := i * chunkSize
			end := start + chunkSize
			if end > len(buffer) {
				end = len(buffer)
			}

			// Process and invert bits in the chunk
			for j := start; j < end; j++ {
				buffer[j] = ^buffer[j]
			}
		}(i)
	}

	wg.Wait()
}
func (f *FileCryptor) GetCryptorInfo() (map[string]interface{}, error) {
	// Open the file

	cryptor_file := f.FileName + ".cryptor"
	if _, err := os.Stat(cryptor_file); os.IsNotExist(err) {
		return nil, fmt.Errorf("cryptor file not found: %w", err)
	}
	file, err := os.Open(cryptor_file) // Open in read-only mode with appropriate permissions
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close() // Ensure file is closed even in case of errors

	// Read the entire file data
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting file info: %w", err)
	}
	fileSize := fileInfo.Size() // Get file size

	// Create buffer for file data
	buffer := make([]byte, fileSize)

	// Read file data into buffer
	n, err := file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Check if all data was read (optional)
	if int64(n) != fileSize {
		return nil, fmt.Errorf("expected to read %d bytes, but read %d", fileSize, n)
	}

	// // Decrypt data (bitwise NOT) directly in the buffer
	// for i := 0; i < len(buffer); i++ {
	// 	buffer[i] = ^buffer[i]
	// }
	// var dict1 map[string]interface{}
	// err = json.Unmarshal(buffer, &dict1)
	// if err != nil {
	// 	return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	// }

	// Convert decrypted data to JSON
	invertBufferParallel(buffer)
	var dict map[string]interface{}
	//buufer nao is holding text like
	// "{\"chunk_size\": 2097152, \"rotate\": 6, \"file-size\": 914630377, \"encondex\": 139702, \"wrap-size\": 6547}ndex\": 139702, \"wrap-size\": 6547}"
	//"{\"chunk_size\": 2097152, \"rotate\": 6, \"file-size\": 914630377, \"encoding\": \"binary\", \"encrypt_block_index\": 139702, \"wrap-size\": 6547}"
	err = json.Unmarshal(buffer, &dict)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Update wrap-size if necessary
	if f.FileSize > 0 && dict["wrap-size"] != nil {
		wrapSize, ok := dict["wrap-size"].(float64)
		if ok && int(wrapSize) > int(f.FileSize) {
			dict["wrap-size"] = f.FileSize
		}
	}

	return dict, nil
}
