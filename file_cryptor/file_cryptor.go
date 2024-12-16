package file_cryptor

import (
	"encoding/json"
	"fmt"
	"io"

	"os"
	"runtime"
	"sync"
	"time"

	cachery "github.com/hashicorp/golang-lru/v2/expirable"
)

type FileCryptorContent struct {
	*os.File
	FileName          string
	FileSize          uint64 // Use int64 for large file sizes
	PathOfFileCryptor string
	CryptInfo         map[string]interface{}
	BufferCache       []byte
	ChunkSize         int
}

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
		if ok && uint64(wrapSize) > uint64(f.FileSize) {
			dict["wrap-size"] = f.FileSize
		}
	}
	if dict["chunk_size"] != nil {
		dict["chunk_size"] = int(dict["chunk_size"].(float64))
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
func decryptContent(dataEncrypt []byte, chunkSize int, firstData byte) <-chan []byte {
	out := make(chan []byte)
	go func() {
		defer close(out)
		pos := 0
		for pos < len(dataEncrypt) {
			buffer := dataEncrypt[pos : pos+chunkSize]
			bufferLen := len(buffer)

			var ret []byte
			lastBit := (firstData & 0x01) << 7

			for i := 0; i < bufferLen-1; i++ {
				rb := (buffer[i] >> 1) | lastBit
				ret = append(ret, rb)
				lastBit = (buffer[i] & 0x01) << 7
			}

			rb := (buffer[bufferLen-1] >> 1) | (buffer[bufferLen-2] & 0x01)
			ret = append(ret, rb)

			out <- ret
			pos += chunkSize
		}
	}()
	return out
}
func (f *FileCryptor) OpenRead() (*FileCryptorContent, error) {
	// Get cryptor info
	cryptorInfo, err := f.GetCryptorInfo()
	if err != nil {
		return nil, err
	}

	// Get file size
	stat, err := os.Stat(f.FileName)
	if err != nil {
		return nil, fmt.Errorf("error getting file info: %w", err)
	}

	fileSize := uint64(stat.Size()) // Get file size

	// Open the file
	file, err := os.Open(f.FileName) // Open in read-only mode with appropriate permissions
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	// Create a new FileCryptorContent instance
	fc := &FileCryptorContent{
		File:              file,
		FileName:          f.FileName,
		FileSize:          fileSize,
		PathOfFileCryptor: f.FileName + ".cryptor",
		CryptInfo:         cryptorInfo,
		BufferCache:       make([]byte, 0),
		ChunkSize:         cryptorInfo["chunk_size"].(int),
	}

	return fc, nil
}
func (f *FileCryptorContent) Seek(offset int64, whence int) (ret int64, err error) {
	return f.File.Seek(offset, whence)

}
func (f *FileCryptorContent) Close() error {
	return f.File.Close()

}
func (f *FileCryptorContent) Read(b []byte) (n int, err error) {
	pos, err := f.File.Seek(0, io.SeekStart)
	if err != nil {
		return 0, fmt.Errorf("error seeking file: %w", err)
	}

	/*
		if not hasattr(fs, "buffer_cache"):
			fs.seek(0)
			encrypt_data = fs.original_read(fs.cryptor["chunk_size"])
			fs.seek(pos)
			decrypt_data = encrypting.decrypt_content(data_encrypt=encrypt_data,
														chunk_size=fs.cryptor["chunk_size"],
														rota=fs.cryptor['rotate'],
														first_data=fs.cryptor['first-data']
														)
			setattr(fs, "buffer_cache", next(decrypt_data))
		if pos + args[0] <= fs.cryptor["chunk_size"]:
			fs.seek(pos + args[0])
			return fs.buffer_cache[pos:pos + args[0]]
		else:
			# ret_fs.seek(ret_fs.cryptor["chunk_size"])
			n = args[0] + pos - fs.cryptor["chunk_size"]
			fs.seek(fs.cryptor["chunk_size"])
			next_data = fs.original_read(n)
			ret_data = fs.buffer_cache[pos:] + next_data
			return ret_data
	*/
	if len(f.BufferCache) == 0 {
		f.Seek(0, io.SeekStart)

		encrypt_data := make([]byte, f.ChunkSize)
		_, err := f.File.Read(encrypt_data)
		if err != nil {
			return 0, fmt.Errorf("error reading file: %w", err)
		}
		decrypt_data := decryptContent(encrypt_data,
			f.ChunkSize,
			byte(f.CryptInfo["first-data"].(float64)))
		f.BufferCache = append(f.BufferCache, <-decrypt_data...)
	}
	if pos+int64(len(b)) <= int64(f.CryptInfo["chunk_size"].(uint64)) {
		f.Seek(pos+int64(len(b)), io.SeekStart)
		//return fs.buffer_cache[pos:pos + args[0]]
		return copy(b, f.BufferCache[pos:pos+int64(len(b))]), nil

	} else {
		// ret_fs.seek(ret_fs.cryptor["chunk_size"])
		var numOfBytes = int64(len(b)) + pos - int64(f.CryptInfo["chunk_size"].(uint64))
		f.Seek(int64(f.CryptInfo["chunk_size"].(uint64)), io.SeekStart)
		next_data := make([]byte, numOfBytes)
		_, err = f.File.Read(next_data)
		if err != nil {
			return 0, fmt.Errorf("error reading file: %w", err)
		}
		ret_data := make([]byte, len(f.BufferCache[pos:]))
		copy(ret_data, f.BufferCache[pos:])
		ret_data = append(ret_data, next_data...)
		return copy(b, ret_data), nil
	}

}
