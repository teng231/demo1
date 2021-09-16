protoc -I/usr/local/include -I. -I$GOPATH/src \
    --gogofaster_out=plugins=grpc:. demo.proto


protoc-go-tags --dir=.