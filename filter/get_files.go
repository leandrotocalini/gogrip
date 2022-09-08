package filter

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func isEnabled(path string) bool {
	if strings.Contains(path, ".git") {
		return false
	}
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
	return contentType == "text/plain; charset=utf-8"
}

func _getFiles(rootPath string, filesInChan chan string) {
	filepath.Walk(rootPath, func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() && isEnabled(path) {
			filesInChan <- path
		}
		return nil
	})
	close(filesInChan)
}

func getFilesEnabledInPath(rootPath string, buffer int) <-chan string {
	filesInChan := make(chan string, buffer)
	go _getFiles(rootPath, filesInChan)
	return filesInChan
}
