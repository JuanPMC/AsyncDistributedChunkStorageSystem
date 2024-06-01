# Makefile for a simple Go project

# Variables
GOCMD := go
GOBUILD := $(GOCMD) build
BINARY_NAME := superPeer

# Default target
all: build

# Build the project
build:
	$(GOBUILD) -o ./bin/$(BINARY_NAME) nodes.go indeingServer.go handleConnection.go communication.go
# Clean the project
clean:
	rm -f ./bin/$(BINARY_NAME)

# PHONY targets ( tu be able to run thefunctions )
.PHONY: all build clean
