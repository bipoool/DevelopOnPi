#!/bin/bash

CONTAINER_ID=$1

# Run the Ubuntu container
sudo docker stop ubuntu_container_$CONTAINER_ID 