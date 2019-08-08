default: clean server client cli

server:
	go build -o ./dist/server ./src/server/main.go

client:
	go build -o ./dist/client ./src/client/main.go

cli:
	go build -o ./dist/cli ./src/cli/main.go

proto3:
	rm ./src/models/*.pb.go
	protoc -I ./src/models ./src/models/*.proto --go_out=plugins=grpc:./src/models
	sed -i '' 's/,omitempty//g' ./src/models/*.pb.go

.PHONY: clean

clean: 
	rm -f dist/*
