package main

import (
	"fmt"
	"net"
	"os"
	"log"
	"encoding/json"
)

// configuration struct
type Configuracion struct {
    Hostname string   `json:"hostname"`
    Port     string      `json:"port"`
    Peers    []string `json:"peers"`
}

// Global variable to hold the configuration
var globalConfig Configuracion

// loadConfig reads configuration from a JSON file
func loadConfig(filePath string) (Configuracion, error) {
    // Open the JSON file
    file, err := os.Open(filePath)
    if err != nil {
        return Configuracion{}, err
    }
    defer file.Close()

    // Decode the JSON data into a Configuracion struct
    var config Configuracion
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return Configuracion{}, err
    }

    // Print the decoded configuration
    decodedConfig, err := json.MarshalIndent(config, "", "    ")
    if err != nil {
        return Configuracion{}, err
    }

    log.Printf("Decoded Configuration:")
    log.Printf(string(decodedConfig))

    return config, nil
}

func main() {
	// setting up logs
	logFile, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Unable to open log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// load the config
	globalConfig, _ = loadConfig("./config.json")

	// Get the port from the command-line argument
	port := globalConfig.Port

	// Listen for incoming connections on the specified port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Server listening on port %s\n", port)

	for {
		// Accept incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(conn)
	}
}
