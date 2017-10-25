# python-17-hw11 #

## Lab configuration ##

processor       : 0
vendor_id       : GenuineIntel
model name      : Intel(R) Xeon(R) CPU E5-2660 v2 @ 2.20GHz
--
processor       : 1
vendor_id       : GenuineIntel
model name      : Intel(R) Xeon(R) CPU E5-2660 v2 @ 2.20GHz
--
processor       : 2
vendor_id       : GenuineIntel
model name      : Intel(R) Xeon(R) CPU E5-2660 v2 @ 2.20GHz
--
processor       : 3
vendor_id       : GenuineIntel
model name      : Intel(R) Xeon(R) CPU E5-2660 v2 @ 2.20GHz

Mem: 16GB

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

2017/10/25 11:46:20 Processing: data/appsinstalled/20170929000200.tsv.gz file.
2017/10/25 11:46:20 Processing: data/appsinstalled/20170929000000.tsv.gz file.
2017/10/25 11:46:20 Processing: data/appsinstalled/20170929000100.tsv.gz file.
2017/10/25 11:50:39 Total lines 3422995 in file data/appsinstalled/20170929000000.tsv.gz.
2017/10/25 11:50:39 Acceptable error rate (0.000000). Successfull load file data/appsinstalled/20170929000000.tsv.gz.
2017/10/25 11:50:40 Total lines 3422026 in file data/appsinstalled/20170929000200.tsv.gz.
2017/10/25 11:50:40 Acceptable error rate (0.000000). Successfull load file data/appsinstalled/20170929000200.tsv.gz.
2017/10/25 11:50:40 Total lines 3424477 in file data/appsinstalled/20170929000100.tsv.gz.
2017/10/25 11:50:40 Acceptable error rate (0.000000). Successfull load file data/appsinstalled/20170929000100.tsv.gz.

real    4m19.688s
user    5m55.220s
sys     3m19.856s

```

Please compare with python version: https://github.com/LiveStalker/python-17/tree/master/hw9
