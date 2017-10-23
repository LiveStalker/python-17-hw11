.PHONY: getpkgs
getpkgs-proto:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go

.PHONY: genproto
genproto: getpkgs-proto
	protoc --go_out=. *.proto