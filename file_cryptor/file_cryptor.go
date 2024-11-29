package file_cryptor

import (
	"encoding/json"
	"fmt"

	"os"
	"runtime"
	"sync"
	"time"

	cachery "github.com/hashicorp/golang-lru/v2/expirable"
)

var cacher = cachery.NewLRU[string, map[string]interface{}](10000, nil, time.Millisecond*10)
var mutex = sync.RWMutex{}

// declare struct including time create cache (datetime) of map[string]interface{} and data with type is map[string]interface{}
// type CacheCryptorInfoItem struct {
// 	time_create time.Time
// 	data        map[string]interface{}
// }
// type CacheCryptorInfo struct {
// 	data map[string]CacheCryptorInfoItem
// 	mu   sync.RWMutex
// }

// decalre a global cache for cryptor info this is a dict key is string and value is CacheCryptorInfo

// declare FileCryptor struct
type FileCryptor struct {
	FileName string
	FileSize uint64 // Use int64 for large file sizes
}

//var cacheCryptorInfo = CacheCryptorInfo{data: make(map[string]CacheCryptorInfoItem)}

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

	// check if crytot_file has in cache_cryptor_info
	// lock cache_cryptor_info for read
	ret, ok := cacher.Get(cryptor_file)
	if ok {
		return ret, nil
	}
	mutex.RLocker().Lock()
	defer mutex.RLocker().Unlock()
	// check if crytot_file has in cache_cryptor_info after lock
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
		wrapSize, ok := dict["wrap-size"].(uint64)
		if ok && int(wrapSize) > int(f.FileSize) {
			dict["wrap-size"] = f.FileSize
		}
	}
	//before add cache_cryptor_info check size of cache_cryptor_info in bytes if bigger than 100MB remove oldest cache until size is less than 100MB
	//check size of cache_cryptor_info in bytes
	// var size int64

	// if size > 100*1024*1024 { //100MB
	// 	//remove oldest cache until size is less than 100MB
	// 	var oldest_key string
	// 	var oldest_time time.Time
	// 	for k, v := range cache_cryptor_info {
	// 		if oldest_time.IsZero() || v.time_create.Before(oldest_time) {
	// 			oldest_key = k
	// 			oldest_time = v.time_create
	// 		}
	// 	}
	// 	delete(cache_cryptor_info, oldest_key)
	// }
	//lock cache_cryptor_info for write

	// cacheCryptorInfo.data[cryptor_file] = CacheCryptorInfoItem{
	// 	time_create: time.Now(),
	// 	data:        dict,
	// }

	fmt.Println("cache_cryptor_info")
	cacher.Add(cryptor_file, dict)
	return dict, nil

}
