#!/bin/bash

# Remove old docker image
docker rmi -f uptonm/uptonm.io:latest

# Build docker container
DOCKER_BUILDKIT=1 docker build -t uptonm/uptonm.io:latest -f Dockerfile .

# List images
# docker images

# Push built container
docker push uptonm/uptonm.io:latest

echo "Done!"