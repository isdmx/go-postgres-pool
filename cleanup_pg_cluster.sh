#!/bin/bash

# Variables
NETWORK_NAME="pg-cluster"
MASTER_NAME="pg-master"
STANDBY_NAME="pg-standby"

# Stop and remove containers
echo "Stopping and removing containers..."
docker stop $MASTER_NAME $STANDBY_NAME
docker rm $MASTER_NAME $STANDBY_NAME

# Remove Docker volumes
echo "Removing Docker volumes..."
docker volume rm pg-master-data pg-standby-data

# Remove Docker network
echo "Removing Docker network..."
docker network rm $NETWORK_NAME

echo "Cleanup complete!"