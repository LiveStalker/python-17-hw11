# python-17-hw11 #

## Bootstrap project (Debian) ##

```bash
sudo apt-get install protobuf-compiler
go get -u github.com/LiveStalker/python-17-hw11
cd $(go env GOPATH)/src/github.com/livestalker/python-17-hw11
make getpkgs
export PATH=$PATH:$(go env GOPATH)/bin
make genproto
```
