#!/bin/bash
set -e

echo "Updating Go modules..."
go get -u ../.
echo "Modules updated successfully."

echo "Tidying up Go modules..."
go mod tidy
echo "Modules tidied successfully."

echo "Formatting Go source files..."
gofmt -s -w ../.
echo "Formatting done."

echo "Building CrunchyCleaner for Windows amd64..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../bin/crunchycleaner.exe ../.
echo "Windows build succeeded."

echo "Building CrunchyCleaner for Linux amd64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../bin/crunchycleaner ../.
echo "Linux amd64 build succeeded."

echo "Building CrunchyCleaner for Linux ARM64..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ../bin/crunchycleaner-arm64 ../.
echo "Linux ARM64 build succeeded."

echo "All builds finished successfully."
