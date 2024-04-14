#!/bin/bash

# Build the main application image
docker build -t codeexec .

# Build the images specified in the dockerfiles directory
for file in dockerfiles/Dockerfile.*; do
    tag=$(echo ${file} | cut -d'.' -f2)
    docker build -t ${tag}-exec -f ${file} .
done
