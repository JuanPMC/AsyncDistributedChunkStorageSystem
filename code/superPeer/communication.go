package main

import (
	"fmt"
	"io"
	"net"
	"crypto/md5"
	"log"
	"strconv"
)

// this is for sending commands to other peers or to the indexing server
func sendCommand(host string,command string) ([]byte, string, error) {
	// log command parameters
	log.Printf(command)
	log.Printf("\n Connecting to: %s \n", host)

	// start a tcp connection
	conn, err := net.Dial("tcp", fmt.Sprintf("%s", host))
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

// function for quering a file in index
func broadcastQuery(peer string, message_id string,sender_id string,ttl int, filename string) ([]byte, string, error){
	// change ttl to string
	ttl_string := strconv.Itoa(ttl)
	// create command
	command := "query "+ message_id +" "+sender_id+" "+ttl_string+" "+ filename
	// send the command and return output
	return sendCommand(peer, command)
}

// send a query hit to the original weak node that requested the file
func queryhit(peer string,message_id string,response string){
	command := "queryhit "+ message_id +" "+ response
	sendCommand(peer, command)
}
