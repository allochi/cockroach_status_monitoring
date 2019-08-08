proto3:
	rm ./src/models/*.pb.go
	protoc -I ./src/models ./src/models/*.proto --go_out=plugins=grpc:./src/models
	# ls ./src/models/*.pb.go | xargs -n1 -IX bash -c "sed -e '/bool/ s/,omitempty//' X > X.tmp && mv X{.tmp,}"
	sed -i '' 's/,omitempty//g' ./src/models/*.pb.go
