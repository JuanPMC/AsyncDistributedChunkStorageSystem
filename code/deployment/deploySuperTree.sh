#!/bin/bash


# check if go is installed
if command -v go &> /dev/null ; then
    echo "Go is installed."
    go version
else
    # Install Go
    echo "Installing Go ..."
    wget -q https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
    rm go1.22.0.linux-amd64.tar.gz

    # Set up Go environment variables
    echo "Setting up Go environment variables..."
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
    echo "export GOPATH=\$HOME/go" >> ~/.bashrc
    echo "export PATH=\$PATH:\$GOPATH/bin" >> ~/.bashrc
    source ~/.bashrc

    # Verify Go installation
    echo "Verifying Go installation..."
    go version
fi

# compile the Weak Peer
cd ../node
make

# compile the Super Peer
cd ../superPeer
make

cd ../deployment


# Define the number of nodes
NUM_NODES=9

# Loop through each node
for ((i=0; i<NUM_NODES; i++)); do
    # Run the superPeer for the current node in the background
    (
        cd "./tree/node$i"
        ../../../superPeer/bin/superPeer
    ) &
    # run all the clients
    (
        cd "./tree/node$i/peer/"
        ../../../../node/bin/node
    ) &
done

# Wait for all background processes to finish
wait
