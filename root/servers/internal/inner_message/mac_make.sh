#!/bin/bash

./protoc-mac --plugin protoc-gen-go=./protoc-gen-go-mac -I=./ --go_out=../inner_message/ ./inner/*.proto