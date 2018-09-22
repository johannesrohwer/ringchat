all: protoc
	go build
run: protoc
	go run *.go
protoc:
	protoc ringchat.proto --go_out=plugins:./grpc/ringchat
clean:
	rm -f ./ringchat
	rm -f ./grpc/ringchat/ringchat.pb.go
