#!/bin/bash

IMAGE_NAME="fetch-api"
CONTAINER_NAME="fetch-api-container"
PORT=8080

function build_image() {
    echo "Building Docker image..."
    docker build -t $IMAGE_NAME .
    if [ $? -ne 0 ]; then
        echo "Failed to build Docker image. Exiting."
        exit 1
    fi
    echo "Docker image built successfully."
}

function run_container() {
    echo "Running Docker container..."
    docker rm -f $CONTAINER_NAME 2>/dev/null

    docker run -d --name $CONTAINER_NAME -p $PORT:$PORT $IMAGE_NAME
    if [ $? -ne 0 ]; then
        echo "Failed to run Docker container. Exiting."
        exit 1
    fi
    echo "Docker container is running on port $PORT."
}

build_image
run_container
