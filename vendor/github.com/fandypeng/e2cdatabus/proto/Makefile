all:clean protos

protos:
	protoc --proto_path=. --gofast_out=plugins=grpc:. *.proto

clean:
	-rm *.pb.go
