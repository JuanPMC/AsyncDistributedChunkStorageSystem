# Makefile for a simple Go project

# Variables
GOCMD := go
GOBUILD := $(GOCMD) build
BINARY_NAME := node

# Default target
all: build

# Build the project
build:
	$(GOBUILD) -o ./bin/$(BINARY_NAME) server.go handle_conn.go main.go monitor_fs.go client.go client_utils.go

# Clean the project
clean:
	rm -f ./bin/$(BINARY_NAME)

# PHONY targets ( tu be able to run thefunctions )
.PHONY: all build clean
