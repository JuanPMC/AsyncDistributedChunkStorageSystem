package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
	"os"
	"log"
)

// Function for querying a file from the server
func queryFileFromServer(host string, port int,filename string){
	idStr := strconv.Itoa(int(time.Now().UnixNano()))
	// Construct the command to query the file
	command := "query 127.0.0.1:8080:" + idStr + " 127.0.0.1:8080 10 " + filename

	// Send the command to the server and retrieve the response
	sendCommand(host, port, command)

}

// Function for sending commands to the server
func sendCommand(host string, port int, command string) ([]byte, string, error) {
	// Establish a TCP connection to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, "", err
	}
	defer conn.Close()

	// Write the command to the server
	_, err = conn.Write([]byte(command))
	if err != nil {
		return nil, "", err
	}

	// Read the response from the server
	var response []byte
	md5Hash := md5.New()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, "", err
		}
		if n == 0 {
			break
		}
		response = append(response, buffer[:n]...)
		md5Hash.Write(buffer[:n])
	}

	md5Hex := fmt.Sprintf("%x", md5Hash.Sum(nil))
	return response, md5Hex, nil
}


func main() {
    // Define the port to listen on
    host := "127.0.0.1"
    port := 9000
	n, _ := strconv.Atoi(os.Args[1])
	// setting up logs
	logFile, err := os.OpenFile("client.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Unable to open log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

    // Repeat the test n times with queries in separate goroutines
    for i := 0; i < n; i++ {
		// Query
		queryFileFromServer(host, port, "shared_dir/pepe3.txt")

		// log the response time
		log.Println(time.Now().UnixNano()/int64(time.Millisecond))

    }

}
