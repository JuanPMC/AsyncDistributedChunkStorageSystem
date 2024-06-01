package main

import (
	"os"
	"path/filepath"
	"sync"
	"time"
	"strconv"
	"log"
)

// a map to store the files (filename, filesize)
var files = make(map[string]struct {
	Size int64
})
// a mux to control the access to the files map
var filesMutex sync.Mutex

// get a list of the files in the directory
func get_current_files(directoryPath string) map[string]struct {Size int64}{
	// Current file listing
	currentFiles := make(map[string]struct {
		Size int64
	})
	// walk all the subdirectoris of the whatched path
	filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		// if its no a directory
		if err == nil && !info.IsDir() {
			// add the file to de list
			currentFiles[path] = struct{ Size int64 }{Size: info.Size()}
		}
		return nil
	})

	// return the list of current files on the server
	return currentFiles
}

// Function to monitor the directory for file changes
func monitorDirectory(directoryPath string) {
	for {
		// check every 0.1 seconds
		time.Sleep(time.Second/10)

		// get list of current files
		currentFiles := get_current_files(directoryPath)

		// Check for changes update the relevant changes in the indexing server
		// lock the files var
		filesMutex.Lock()
		// loop the files
		for file := range files {
			// if file dosent exist in current files
			if _, exists := currentFiles[file]; !exists {
				// the file has been deleted
				log.Printf("File deleted: %s (Size: %d bytes)\n", file, files[file].Size)
				// notify the index server
				unindex_file(globalConfig.Hostname,strconv.Itoa(globalConfig.Port), file)
			}
		}
		// loop current files
		for file := range currentFiles {
			// if current file is not in files
			if _, exists := files[file]; !exists {
				// file was added
				log.Printf("File added: %s (Size: %d bytes)\n", file, currentFiles[file].Size)
				// add file to indexing server
				index_file(globalConfig.Hostname,strconv.Itoa(globalConfig.Port), file, strconv.FormatInt(currentFiles[file].Size,10))
			}
		}

		// Update the files map
		files = currentFiles
		// unlock mutex
		filesMutex.Unlock()
	}
}
