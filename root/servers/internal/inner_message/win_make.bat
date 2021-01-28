protoc --plugin protoc-gen-go=./protoc-gen-go-win.exe -I=./ --go_out=../inner_message/ ./inner/*.proto

pause