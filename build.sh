#!/bin/bash

rm -rf ./dist/*
echo Building for Windows
GOOS=windows go build -o dist/signalred.exe
echo Building for Linux
GOOS=linux go build -o dist/signalred-linux
echo Building for Mac
GOOS=darwin go build -o dist/signalred-darwin
echo Done