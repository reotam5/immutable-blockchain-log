#!/bin/bash

BACKUP_DIR="config_backup"

# check if bin directory exists. If it is, exit the script.
if [ -d "bin" ]; then
  echo "Binaries already installed, exiting..."
  exit 0
fi

# Remove existing backup if it exists
if [ -d "$BACKUP_DIR" ]; then
  echo "Removing existing backup directory..."
  rm -rf "$BACKUP_DIR"
fi

if ! mkdir -p "$BACKUP_DIR"; then
  echo "Error: Failed to create backup directory $BACKUP_DIR"
  return 1
fi

# Backup config directory if it exists
echo "Backing up config/ directory..."
if ! cp -r config "$BACKUP_DIR/"; then
  echo "Error: Failed to backup config/ directory"
  return 1
fi


echo "Downloading official Hyperledger Fabric install script..."
# Download the official install script
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh

if [ $? -eq 0 ]; then
  echo "Running official install script to download binaries..."
  # Install only binaries (b) with specified versions
  ./install-fabric.sh binary
  
  if [ $? -eq 0 ]; then
    # Clean up the install script
    rm -f install-fabric.sh
    
    # Replace config directory with the backup
    if [ -d "$BACKUP_DIR/config" ]; then
      rm -rf config
      if ! cp -r "$BACKUP_DIR/config" .; then
        return 1
      fi
    fi
    rm -rf "$BACKUP_DIR"
    
    echo "Installation complete!" 
  else
    echo "Error: Binary installation failed"
    rm -f install-fabric.sh
    exit 1
  fi
else
  echo "Error: Failed to download official install script"
  exit 1
fi

