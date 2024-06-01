package main

import (
	"net"
	"io/ioutil"
	"strings"
	"log"
	"time"
	"crypto/md5"
	"encoding/hex"
)


// generates hexed md5 string
func calculateMD5(filename string) (string) {
	// reads file
	content, _ := ioutil.ReadFile(globalConfig.WatchPath+"/"+filename)
	// generates md5 sum
	sum := md5.Sum(content)
	return hex.EncodeToString(sum[:])
}

func download_file(filename string) ([]byte) {
	// Read the entire file into a byte slice
	content, err := ioutil.ReadFile(globalConfig.WatchPath+"/"+filename)
	if err != nil {
		return []byte("ERROR: File not fund")
	}

	return content
}

// this function processes the commands form the client
func process_comand(bytesRead int,buffer []byte) ([]byte, string){
	// reeds the command and stores it in orden
	orden := strings.Fields(string(buffer[:bytesRead]))
	// Defines the default output and extra logs
	output := []byte("orden dosent exist\n")
	extra_logs := ""

	// command download
	if orden[0] == "download"{
		// if filename is included
		if len(orden) >= 2{
			// call the download funbction
			output = download_file(orden[1])
			// put the filename in the extralogs
			extra_logs = "Filename: " + orden[1]
		}
	}
	// command md5
	if orden[0] == "md5"{
		// if file is specified
		if len(orden) >= 2{
			// calculate md5 callign the apropiate function
			output = []byte(calculateMD5(orden[1]))
			log.Printf("Calculado el hash: %x",string(output))
		}
	}
	// command queryhit
	if orden[0] == "queryhit"{
		// if file is specified
		if len(orden) >= 3{
			// add file to the list
			update_queryhits(orden[1],strings.Join(orden[2:], ""))
			// send a response
			output = []byte("ok")
			log.Printf("Got a query hit: ID=" + orden[1] +" RESPONSE:"+orden[2])
		}
	}
	// return the output and the extralogs
	return output, extra_logs
}

// No loop to make it rest-full
func handleConnection(conn net.Conn) {
	// close the connecton at the end of the function
	defer conn.Close()
	var extra_logs string

	// Log the client connection
	clientAddr := conn.RemoteAddr()
	log.Printf("Client connected: %s\n", clientAddr)

	buffer := make([]byte, 1024)

	// Read data from the connection
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		log.Println("Error reading:", err)
		return
	}

	output, extra_logs := process_comand(bytesRead,buffer)

	// Return the output of the processed request
	startTime := time.Now()
	_, err = conn.Write(output)
	if err != nil {
		log.Println("Error writing:", err)
		return
	}
	// log the ammount of bytes sent
	log.Printf("Bytes: %d Client: %s Transfer_Time: %s %s\n", len(output),clientAddr, time.Since(startTime),extra_logs)

}
