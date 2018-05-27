gen-proto:
	@protoc -I api/ --go_out=plugins=grpc:api api.proto
