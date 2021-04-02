#!/bin/bash

go build -o bloggulus main.go
./bloggulus -syncblogs
