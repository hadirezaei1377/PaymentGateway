#!/bin/bash

# Stop the running application (if any)
sudo systemctl stop your-application

# Pull the latest changes from the repository
git pull

# Build the application
go build main.go

# Restart the application
sudo systemctl start your-application
