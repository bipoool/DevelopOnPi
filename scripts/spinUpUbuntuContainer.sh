#!/bin/bash

# Set your SSH public key here
SSH_PORT=$1
CPU_COUNT=$2
MEMORY="${3}g"
SSH_PUBLIC_KEY=$4
CONTAINER_ID=$5

# Run the Ubuntu container
sudo docker run -d --name ubuntu_container_$CONTAINER_ID \
    -p $SSH_PORT:22 \
    --cpus=$CPU_COUNT \
    --memory=$MEMORY \
    ubuntu:20.04 sleep infinity

# Wait for the container to start and be ready
sleep 5

# Install openssh-server inside the container
# sudo docker exec ubuntu_container_$CONTAINER_ID bash -c "apt-get update -y && DEBIAN_FRONTEND=noninteractive apt-get install -y openssh-server"

# Create the necessary directories for SSH if they don't exist
sudo docker exec ubuntu_container_$CONTAINER_ID bash -c "mkdir -p /root/.ssh"

# Add the public key to /root/.ssh/authorized_keys
sudo docker exec ubuntu_container_$CONTAINER_ID bash -c "echo '$SSH_PUBLIC_KEY' >> /root/.ssh/authorized_keys"

# Set the correct permissions for the SSH folder and authorized_keys file
sudo docker exec ubuntu_container_$CONTAINER_ID bash -c "chmod 700 /root/.ssh && chmod 600 /root/.ssh/authorized_keys"

# Start the SSH service inside the container
sudo docker exec ubuntu_container_$CONTAINER_ID bash -c "service ssh start"

# Print the container's IP and the SSH command to connect
echo "SSH setup complete! You can now SSH into the container using:"
echo "ssh root@localhost -p $SSH_PORT"
