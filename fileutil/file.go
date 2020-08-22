package fileutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	max_file_size_in_bytes = 4 * (2 << 20)
	SIZE_TYPE_LARGE        = "LARGE"
	SIZE_TYPE_SMALL        = "SMALL"
	default_chunk_size     = (60 * 1024 * 1024)
)

type FileInfo struct {
	FileData *os.File
	SizeType string
}

func GetDefaultChunkSize() int64 {
	return default_chunk_size
}
func GetAllUploadItemsFrmSource(sourcePath string) (map[string]FileInfo, error) {
	fileMap := make(map[string]FileInfo)
	err := filepath.Walk(sourcePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				//Create FileInfo object
				fileInfo := FileInfo{
					SizeType: SIZE_TYPE_SMALL,
				}
				//If file size is greater than 4 mb return error
				//for now until there is a support for Large file upload.
				if info.Size() > max_file_size_in_bytes {
					fileInfo.SizeType = SIZE_TYPE_LARGE
					//return fmt.Errorf("File %s size  %d > 4Mb is not allowed for simple Restore", info.Name(), info.Size())
				}
				fileItem, err := os.Open(path)
				if err != nil {
					return err
				}
				//parentDir := filepath.Dir(path)
				//fmt.Println(parentDir)
				fileInfo.FileData = fileItem
				fileMap[path] = fileInfo
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return fileMap, nil
}

//Returns the file in parts based on the provided offset
func GetFilePartInBytes(filePath string, startingOffset int64, isLastChunk bool) ([]byte, error) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		return nil, err
	}
	var buffer []byte
	if isLastChunk {
		lastChunkSize, err := GetLatsChunkSizeInBytes(filePath)
		if err != nil {
			return nil, err
		}
		buffer = make([]byte, lastChunkSize)
	} else {
		buffer = make([]byte, default_chunk_size)
	}

	_, err = file.ReadAt(buffer, startingOffset)
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("readAt: %v", err)
		}
	}
	return buffer, nil
}

//Returns the start offset chunk list based on the file size
func GetFileOffsetStash(filePath string) ([]int64, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	//Get file size
	size, err := GetFileSize(filePath)
	if err != nil {
		return nil, err
	}

	//Get the max offset length to calculate the chunks
	offsetMax := size - 1

	//Based on the offsetMax generate the start offset list
	var i int64
	offsetLst := make([]int64, 0)
	for i = 0; i <= offsetMax; i = i + default_chunk_size {
		offsetLst = append(offsetLst, i)
	}
	return offsetLst, nil
}

//Get file size
func GetFileSize(filePath string) (int64, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return -1, err
	}
	//Get the file size in bytes
	fi, err := file.Stat()
	if err != nil {
		return -1, err
	}

	fSize := fi.Size()
	return fSize, nil
}

//return the last chunk size based on modulus with the file size with the default chunk size.
//If the modulus is 0 then that means the file size is properly divisible by default chunk
//size and so return the default chunk size in that case.
func GetLatsChunkSizeInBytes(filePath string) (int64, error) {
	fSize, err := GetFileSize(filePath)
	if err != nil {
		return -1, err
	}
	var lastChunkSize int64
	chunkSize := (fSize % default_chunk_size)
	if chunkSize == 0 {
		lastChunkSize = default_chunk_size
	} else {
		lastChunkSize = chunkSize - 1
	}
	return lastChunkSize, nil
}
func GetAlternateRootFolder() string {
	dt := time.Now()
	//restore_YYYYMMDD_hhmmssff
	return fmt.Sprintf("restore_%s", dt.Format("20060102_15040535"))
}

// Read the file content from File handle
func ReadFile(file *os.File) ([]byte, error) {
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
