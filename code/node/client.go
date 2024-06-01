package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"io/ioutil"
	"strconv"
	"log"
)

// get list of files
func getFilesList() {
	// call list_files to get the list of files
	response, _, _ := list_files()
	// show the list of files
	fmt.Println("Files List:")
	fmt.Println(string(response))
}

// to download files
func download(fileName string) bool {
	// log the begin download time
	startTime := time.Now()

	// get peers that contain the file
	servers := query_file(fileName)
	// pick out the first one
	server := strings.Split(string(servers), "|")[0]

	// get the ip of the peer
	ip := strings.Split( server, ":")[0]
	// get port of peer
	port, _:= strconv.Atoi(strings.Split( server, ":")[1])

	// get the actual filename ( not the path )
	filepath := strings.Split(fileName,"/")
	fileName = filepath[len(filepath)-1]

	// request the md5 hash of the file to the peer
	// create command
	md5Command := "md5 " + fileName
	// send command to peer
	sumReal, _, _:= sendCommand(ip,port,md5Command)
	// get the resulting hash
	sumRealStr := string(sumReal)

	// send the download commmand to the peer
	downloadCommand := "download " + fileName
	// store the response file and its md5hash
	response, sum, _ := sendCommand(ip,port,downloadCommand)

	log.Printf(": sum: %s sum_real: %s\n", sum, sumRealStr)

	// compair the hashes
	if strings.HasPrefix(string(response), "Error") || sum != sumRealStr {
		log.Println("Download Error")
		return false
	}

	// write the file to the Download_dir
	// define file path
	filePath := globalConfig.WatchPath + "/" + fileName
	// write create the file
	err := ioutil.WriteFile(filePath, response, 0644)
	// check for an error
	if err != nil {
		log.Println("Error writing file:", err)
		return false
	}

	// log the finish time
	endTime := time.Now()
	// compute the download time ( important for evaluations )
	downloadTime := endTime.Sub(startTime).Seconds()
	// log the time in the log file
	log.Printf("Downloaded %s in %.2f seconds.\n", fileName, downloadTime)

	return true
}

// Donwload the files in series
func downloadSeries(filesToDownload []string) {
	// get the bigging time
	startTime := time.Now()

	// loop all requested files
	for _, fileName := range filesToDownload {
		// call the download function for all of them
		if !download(fileName) {
			fmt.Println("Error in file integrity, download stopped")
			return
		}
	}

	// compute the total time for all files
	totalDownloadTime := time.Since(startTime).Seconds()
	// print the information on the download time
	log.Println("\nSummary series:")
	log.Printf("Downloaded Files: %s\n", strings.Join(filesToDownload, ", "))
	log.Printf("Total Download Time: %.2f seconds\n", totalDownloadTime)
}

// download files in parallel
func downloadInParallel(filesToDownload []string) {
	// get the bigging time
	startTime := time.Now()

	// Create a WaitGroup to wait for all files to be downloaded
	var wg sync.WaitGroup

	// loop files
	for _, fileName := range filesToDownload {
		// increment the wait group
		wg.Add(1)
		// create a thrad of an anonymous function
		go func(file string) {
			// defer the signal to the wg
			defer wg.Done()
			// call download to download the file
			download(file)
		}(fileName)
	}

	// Wait for all files to be downloaded
	wg.Wait()

	// compute thte total time
	totalDownloadTime := time.Since(startTime).Seconds()
	// print the information on the download time
	log.Println("\nSummary parallel:")
	log.Printf("Downloaded Files: %s\n", strings.Join(filesToDownload, ", "))
	log.Printf("Total Download Time: %.2f seconds\n", totalDownloadTime)
}

// the main client function
func client_main() {
	// load the config
	globalConfig, _ = loadConfig("./config.json")
	// create a reader for stdin ( to get user input )
	reader := bufio.NewReader(os.Stdin)
	// Create a loop for the interface
	for {
		// show menu
		fmt.Println("\nOptions:")
		fmt.Println("1. Get Files List")
		fmt.Println("2. Download File(s) in series")
		fmt.Println("3. Download File(s) in parallel")
		fmt.Println("4. Exit")

		fmt.Print("Enter your choice (1/2/3/4): ")
		// store the choice
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		// Created a switch case for all the input parameters
		switch choice {
		case "1":
			fmt.Println("--------------------------------")
			// call to get the files listed
			getFilesList()
			fmt.Println("--------------------------------")
		case "2":
			// created sub menu for inputing multiple files
			fmt.Print("Enter the file names to download (comma-separated): ")
			// process the listed files
			filesToDownloadStr, _ := reader.ReadString('\n')
			filesToDownload := strings.Split(strings.TrimSpace(filesToDownloadStr), ",")
			fmt.Println("--------------------------------")
			// call download in series
			downloadSeries(filesToDownload)
			fmt.Println("--------------------------------")
		case "3":
			// created sub menu for inputing multiple files
			fmt.Print("Enter the file names to download (comma-separated): ")
			// process the listed files
			filesToDownloadStr, _ := reader.ReadString('\n')
			filesToDownload := strings.Split(strings.TrimSpace(filesToDownloadStr), ",")
			fmt.Println("--------------------------------")
			// call download in parallel
			downloadInParallel(filesToDownload)
			fmt.Println("--------------------------------")
		case "4":
			// exit the program
			fmt.Println("Exiting the program.")
			return
		default:
			// print out htat no coice was selected
			fmt.Println("--------------------------------")
			fmt.Println("Invalid choice. Please enter 1, 2, 3, or 4")
			fmt.Println("--------------------------------")
		}
	}
}