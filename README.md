# python-17-hw11 #

## Bootstrap project (Debian) ##

```
sudo apt-get install protobuf-compiler
go get -u github.com/LiveStalker/python-17-hw11
cd $(go env GOPATH)/src/github.com/livestalker/python-17-hw11
make getpkgs
export PATH=$PATH:$(go env GOPATH)/bin
make genproto
go install
```

Plase download test data.

## Run ##

```
make run_memcached
ls -l ./data/appsinstalled/
total 1555076
-rw-r--r-- 1 aleksio aleksio 530795343 Oct 24 15:18 20170929000000.tsv.gz
-rw-r--r-- 1 aleksio aleksio 530811684 Oct 24 15:30 20170929000100.tsv.gz
-rw-r--r-- 1 aleksio aleksio 530787274 Oct 24 15:28 20170929000200.tsv.gz
time python-17-hw11 --pattern './data/appsinstalled/*.tsv.gz'
```
