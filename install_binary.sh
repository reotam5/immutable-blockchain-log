#!/bin/bash

# check if bin directory exists. If it is, exit the script.
if [ -d "bin" ]; then
  exit 0
fi

echo "Downloading official Hyperledger Fabric install script..."
# Download the official install script
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh

if [ $? -eq 0 ]; then
  echo "Running official install script to download binaries..."
  # Install only binaries (b) with specified versions
  ./install-fabric.sh binary
  
  # Clean up the install script
  rm -f install-fabric.sh
  
  echo "Installation complete!" 
else
  echo "Failed to download official install script"
  exit 1
fi

