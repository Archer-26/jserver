protoc --plugin protoc-gen-go=./protoc-gen-go-win.exe -I=./ --go_out=./ ./protofile/*.proto
pause