package main

import (
	"os"
	"flag"
	"fmt"
	"runtime"
	"log"
	"sort"
	"bufio"
	"path/filepath"
	"compress/gzip"
	"github.com/bradfitz/gomemcache/memcache"
	//"strings"
	//"time"
	//"./appsinstalled"
	//"github.com/golang/protobuf/proto"
)

func main() {
    memc := make(map[string]string)
	pattern := flag.String("pattern", "", "Files pattern.")
	idfaOpt := flag.String("idfa", "127.0.0.1:33013", "Memcache for idfa.")
	gaidOpt := flag.String("gaid", "127.0.0.1:33014", "Memcache for gaid.")
	adidOpt := flag.String("adid", "127.0.0.1:33015", "Memcache for adid.")
	dvidOpt := flag.String("dvid", "127.0.0.1:33016", "Memcache for dvid.")
	workers := flag.Int("workers", runtime.NumCPU(), "Count of forkers.")
	fmt.Println(*pattern, *idfaOpt, *gaidOpt, *adidOpt, *dvidOpt, *workers)
	flag.Parse()
	if *pattern == "" {
		log.Fatal("Pattern not found in arguments")
	}
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	} else {
        memc["idfa"] = *idfaOpt
        memc["gaid"] = *gaidOpt
        memc["adid"] = *adidOpt
        memc["dvid"] = *dvidOpt
	}
	runtime.GOMAXPROCS(*workers)
	start(pattern, &memc)
}

func start(pattern *string, memc *map[string]string) {
	mClients := make(map[string]*memcache.Client)
	for key, value := range *memc {
		mClients[key] = memcache.New(value)
	}
	files, err := filepath.Glob(*pattern); if err != nil {
		log.Fatal(err)
	}
	sort.Strings(files)
	for _, f := range files {
		log.Printf("Processing: %s file.", f)
		fh, err := os.Open(f); if err != nil {
			log.Printf("File: %s, error: %s", f, err)
			continue
		}
		gz, err := gzip.NewReader(fh); if err != nil {
			log.Println(err)
			fh.Close()
			continue
		}
		scanner := bufio.NewScanner(gz)

		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		gz.Close()
		fh.Close()
	}
}
