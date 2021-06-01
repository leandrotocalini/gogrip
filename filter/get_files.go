package filter

import (
	"net/http"
	"os"
	"path/filepath"
)

func isText(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	buffer := make([]byte, 16)
	_, err = file.Read(buffer)
	if err != nil {
		return false
	}
	file.Seek(0, 0)
	contentType := http.DetectContentType(buffer)
	if contentType == "text/plain; charset=utf-8" {
		return true
	}
	return false
}

func _getFiles(rootPath string, filesInChan chan string) {
	filepath.Walk(rootPath, func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			if isText(path) {
				filesInChan <- path
			}
		}
		return nil
	})
	close(filesInChan)
}

func getFiles(rootPath string, buffer int) <-chan string {
	filesInChan := make(chan string, buffer)
	go _getFiles(rootPath, filesInChan)
	return filesInChan
}
