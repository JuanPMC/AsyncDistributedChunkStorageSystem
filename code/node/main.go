package main

import(
	"os"
	"encoding/json"
	"log"
	"strconv"
	"sync"
)

// the configuration object
// all variables have been mapped to json
type Configuracion struct {
    Hostname   string   `json:"hostname"`
    Port int `json:"port"`
	WatchPath string `json:"watchpath"`
	Indexing_server string `json:"indexing_server"`
	Indexing_port int `json:"indexing_port"`
	Download_dir string `json:"download_dir"`
}

// Global variable to hold the configuration
var globalConfig Configuracion
// Global variables to handle the queryhits
var (
	// map messafeid to responses
	queryhits = make(map[string]string)
	currentid = 0
	// query_mu to control acces to the queryhits
	query_mu	sync.Mutex
)


// load the configuration from the json file
func loadConfig(filePath string) (Configuracion, error) {
	// Open the JSON file
	file, err := os.Open(filePath)
	// return error if an error ocurrs when opening
	if err != nil {
		return Configuracion{}, err
	}
	// defer the closing of the file
	defer file.Close()

	// Decode the JSON data into a Configuracion struct
	var config Configuracion
	// create decoder
	decoder := json.NewDecoder(file)
	// decode the json into the config object
	if err := decoder.Decode(&config); err != nil {
		return Configuracion{}, err
	}

	// Print the decoded configuration
	decodedConfig, err := json.MarshalIndent(config, "", "    ")
	// return if any error unmarshaling
	if err != nil {
		return Configuracion{}, err
	}

	log.Printf("Decoded Configuration:")
	// print the configuration in the log
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

	// start server on a thrad
	go server_main()
	// start client
	client_main()
	// unregister server before exiting
	unregister(globalConfig.Hostname,strconv.Itoa(globalConfig.Port))
	os.Exit(1)
}