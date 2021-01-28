#!/bin/bash

./protoc-mac --plugin protoc-gen-go=./protoc-gen-go-mac -I=./ --go_out=./ ./protofile/*.proto