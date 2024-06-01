package main

import (
	"fmt"
	"io"
	"net"
	"crypto/md5"
	"log"
	"time"
	"strconv"
)

// this is for sending commands to other peers or to the indexing server
func sendCommand(host string, port int,command string) ([]byte, string, error) {
	// log command parameters
	log.Printf(command)
	log.Printf("\n Connecting to: %s:%d \n", host, port)

	// start a tcp connection
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	// check for errors in the connexion
	if err != nil {
		return nil, "", err
	}
	// defere the conexion close
	defer conn.Close()

	// write the command to the server
	_, err = conn.Write([]byte(command))
	if err != nil {
		return nil, "", err
	}

	// declare response variable
	var response []byte
	// declare md5hash
	md5Hash := md5.New()
	// declare buffer
	buffer := make([]byte, 1024)

	// read response
	for {
		// get a buffer size chunck
		n, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			// return the error
			return nil, "", err
		}
		// if no more bytes were added
		if n == 0 {
			break
		}
		// append to response
		response = append(response, buffer[:n]...)
		// compute hash
		md5Hash.Write(buffer[:n])
	}

	// create hash as string
	md5Hex := fmt.Sprintf("%x", md5Hash.Sum(nil))
	// return response and hash
	return response, md5Hex, nil
}

// function for registering with indeing server
func register(host string, port string){
	// create command
	command := "register " + host + " " + port + "\n"
	// send the command and return output
	sendCommand(globalConfig.Indexing_server, (globalConfig.Indexing_port), command)
}
// function for registering with indeing server

func unregister(host string, port string){
	// create command
	command := "unregister " + host + ":" + port + "\n"
	// send the command and return output
	sendCommand(globalConfig.Indexing_server, (globalConfig.Indexing_port), command)
}

// function for updating a file in index
func index_file(host string, port string, filename string, size string){
	// create command
	command := "add " + host + ":" + port + " " + filename + " " + size + "\n"
	// send the command and return output
	sendCommand(globalConfig.Indexing_server, (globalConfig.Indexing_port), command)
}
// function for removing a file in index
func unindex_file(host string, port string, filename string){
	// create command
	command := "rm " + host + ":" + port + " " + filename + "\n"
	// send the command and return output
	sendCommand(globalConfig.Indexing_server, (globalConfig.Indexing_port), command)
}

// function to update query hits
func update_queryhits(msgid string, text string){
	query_mu.Lock()
	defer query_mu.Unlock()
	// if exist append
	if _, exists := queryhits[msgid]; exists {
		queryhits[msgid] += text + "|"
	}
	queryhits[msgid] = text + "|"
}

// function for quering a file in index
func query_file(filename string) (string){

	// compute the msg id
	msgid := globalConfig.Hostname +":" + strconv.Itoa(globalConfig.Port) +":" + strconv.Itoa(currentid)

	// setting ttl as 10
	command := "query "+ msgid +" " +globalConfig.Hostname +":" + strconv.Itoa(globalConfig.Port)+" 10 "+ filename

	// incrementing the current id
	currentid +=1

	// send the query command
	sendCommand(globalConfig.Indexing_server, (globalConfig.Indexing_port), command)
	// wait for respionse
	time.Sleep(2000 * time.Millisecond)

	// send the command and return output
	return queryhits[msgid]
}
// function for getting list of files from index
func list_files()([]byte, string, error){
	// create command
	command := "state"
	// send the command and return output
	return sendCommand(globalConfig.Indexing_server, (globalConfig.Indexing_port), command)
}