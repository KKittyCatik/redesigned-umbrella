#!/bin/bash

# This script sets up the development environment for the Go project.

# Update and install necessary packages
sudo apt-get update
sudo apt-get install -y golang-go

# Set up Go workspace
mkdir -p $HOME/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Install necessary Go tools
go get -u github.com/gorilla/mux
go get -u github.com/stretchr/testify

echo "Development environment setup complete."