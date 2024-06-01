package main

import (
	"net"
	"os"
	"log"
	"strconv"
	"os/signal"
)

func server_main() {
	// Get the port from the command-line argument
	port := strconv.Itoa(globalConfig.Port)

	// Listen for incoming connections on the specified port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Printf("Error listening:", err)
		return
	}
	defer listener.Close()
	log.Printf("Server listening on port %s\n", port)

	// register with index
	register(globalConfig.Hostname,strconv.Itoa(globalConfig.Port))
	defer unregister(globalConfig.Hostname,strconv.Itoa(globalConfig.Port))

	// Make shure the node unregisters
	c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt)
    go func() {
        select {
        case sig := <-c:
            log.Printf("Got %s signal. Aborting...\n", sig)
			unregister(globalConfig.Hostname,strconv.Itoa(globalConfig.Port))
            os.Exit(1)
        }
    }()


	// start nmonitoring the filesystem
	go monitorDirectory(globalConfig.WatchPath)

	for {
		// Accept incoming connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(conn)
	}
}
