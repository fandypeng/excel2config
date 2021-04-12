all:clean protos

protos:
	protoc --go_out=plugins=grpc:. *.proto

clean:
	-rm *.pb.go
