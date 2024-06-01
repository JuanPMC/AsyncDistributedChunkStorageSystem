package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"log"
)

func main() {
	// setting up logs
	logFile, err := os.OpenFile("server.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Unable to open log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Listen for incoming connections on port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()

	fmt.Println("Listening on port 8080")

	// Accept connections in a loop
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return
		}

		// Log the milliseconds when the connection is opened to the file
		log.Println(time.Now().UnixNano()/int64(time.Millisecond))
		if err != nil {
			fmt.Println("Error writing to file:", err.Error())
			return
		}

		// Close the connection
		conn.Close()
		fmt.Println("Aceptao")

	}
}
