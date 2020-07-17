#!/bin/bash

# Build.
echo "Starting build..."

if [ ! -d "dist" ]; then
    if ! mkdir "dist"; then
        echo "Failed to create dist directory, exiting..."; exit 1
    fi
fi

echo "Building binary"
if ! go build -o main cmd/popper/main.go; then
    echo "Failed to build binary, exiting..."; exit 1
fi

echo "Zipping binary"
if ! zip dist/main.zip main > /dev/null; then
    echo "Failed to zip binary, exiting..."; exit 1
fi

echo "Removing unzipped binary"
if ! rm main; then
    echo "Failed to remove unzipped binary, build still successfully completed"
fi

echo "Build completed!"

# Deploy.
echo "Starting deployment..."

functionName=

case "$1" in
    dev)
        functionName=MailService_Dev
        ;;
    prod)
        functionName=MailService_Prod
        ;;
    *)
        echo "Failed to deploy zipped binary, invalid deployment parameter: '$1', exiting..."; exit 1
esac

if [ -n "$2" ]; then
    export AWS_PROFILE=$2
fi

if ! aws lambda update-function-code --function-name $functionName --zip-file fileb://dist/main.zip; then
    echo "Failed to deploy zipped binary"; exit 1
fi

echo "Deployment completed!"

