package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	max_file_size_in_bytes = 4 * (2 << 20)
)

func GetAllUploadItemsFrmSource(sourcePath string) (map[string]*os.File, error) {
	fileMap := make(map[string]*os.File)
	err := filepath.Walk(sourcePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				//If file size is greater than 4 mb return error
				if info.Size() > max_file_size_in_bytes {
					return fmt.Errorf("File %s size  %d > 4Mb is not allowed for simple Restore", info.Name(), info.Size())
				}
				fileItem, err := os.Open(path)
				if err != nil {
					return err
				}
				//parentDir := filepath.Dir(path)
				//fmt.Println(parentDir)
				fileMap[path] = fileItem
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return fileMap, nil
}

func GetAlternateRootFolder() string {
	dt := time.Now()

	//restore_YYYYMMDD_hhmmssff
	return fmt.Sprintf("restore_%s", dt.Format("20060102_15040535"))
}
