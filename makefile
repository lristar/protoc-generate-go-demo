init:
	protoc -I ./protos --go_out ./protos --go_opt=paths=source_relative --go-demo_out ./protos --go-demo_opt=paths=source_relative protos/hello.proto