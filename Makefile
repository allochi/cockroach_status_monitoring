proto3:
	protoc -I ./src/models ./src/models/*.proto --go_out=plugins=grpc:./src/models
