package main

import (
	"fmt"
	"sync"
	"log"
)

// Node represents a registered node
type Node struct {
	IP   string // ip
	Port string // port
	ID   string // id of the node ( to identifie node in map )
	fileMap map[string]int // to store the files and the amount of bytes
	mutex   sync.Mutex // mutex to control the updating of files
}

// Add a file to the list
func (n *Node) addFile(filename string, filesize int) {
	// lock the mutex
	n.mutex.Lock()
	defer n.mutex.Unlock()
	// add filesize mapped to filename
	n.fileMap[filename] = filesize
	log.Printf("Node %s added  %s: %d bytes\n", n.ID,filename, filesize)
}

// Delete a file
func (n *Node) deleteFile(filename string) {
	// lock the mutex
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// the lete the entry with the filename as key
	delete(n.fileMap, filename)

	log.Printf("Node %s removed  %s\n", n.ID,filename)
}

// NewNode is a constructor function for creating a Node with a computed ID
func NewNode(IP string, Port string) *Node {
	// set the apropiate variables
	node := &Node{
		IP: IP,
		Port: Port,
		ID: fmt.Sprintf("%s:%s", IP, Port),
		fileMap: make(map[string]int),
	}
	// return the created node
	return node
}

func (n *Node) printInfo() string {
	// initialize the info with ip port and id
	info := fmt.Sprintf("Node ID: %s, IP: %s, Port: %s\n", n.ID, n.IP, n.Port)
	info += "Files:\n"
	// loop files
	for filename, filesize := range n.fileMap {
		// append file info
		info += fmt.Sprintf("  %s: %d bytes\n", filename, filesize)
	}
	// return all the info
	return info
}