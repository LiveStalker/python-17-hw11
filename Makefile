BIND_ADDRESS=127.0.0.1

# Get all needed packages
.PHONY: getpkgs-proto
getpkgs: getpkgs-proto
	go get -u github.com/bradfitz/gomemcache/memcache

# Get packages for protobuf
.PHONY: getpkgs
getpkgs-proto:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go

# Compile protocol buffer description files
.PHONY: genproto
genproto: getpkgs-proto
	-mkdir appsinstalled 
	protoc --go_out=./appsinstalled *.proto

.PHONY: run_memcached
run_memcached:
	@memcached -l $(BIND_ADDRESS) -m 1024 -p 33013 &
	@memcached -l $(BIND_ADDRESS) -m 1024 -p 33014 &
	@memcached -l $(BIND_ADDRESS) -m 1024 -p 33015 &
	@memcached -l $(BIND_ADDRESS) -m 1024 -p 33016 &
	@ps -C memcached

.PHONY: stop_memcaced
stop_memcaced:
	@-kill -s 9 `ps aux | grep [m]emcached | tr -s ' ' | cut -f2 -d' '`
