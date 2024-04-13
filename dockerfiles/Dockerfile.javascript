# Use an official Node.js runtime as a parent image
FROM node:19-alpine

# Set the working directory
WORKDIR /app

# Copy the local code to the container
COPY . /app

# Set the entry point to node so we can run JavaScript files
ENTRYPOINT ["node"]
