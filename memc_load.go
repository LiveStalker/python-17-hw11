package main

import (
	"flag"
	"fmt"
	//"log"
	"runtime"
	//"strings"
	//"time"
)

func main() {
    var memc map[string]string
	pattern := flag.String("pattern", "./data/appsinstalled/*.tsv.gz", "Files pattern.")
	idfaOpt := flag.String("idfa", "127.0.0.1:33013", "Memcache for idfa.")
	gaidOpt := flag.String("gaid", "127.0.0.1:33014", "Memcache for gaid.")
	adidOpt := flag.String("adid", "127.0.0.1:33015", "Memcache for adid.")
	dvidOpt := flag.String("dvid", "127.0.0.1:33016", "Memcache for dvid.")
	workers := flag.Int("workers", runtime.NumCPU(), "Count of forkers.")
	fmt.Println(*pattern, *idfaOpt, *gaidOpt, *adidOpt, *dvidOpt, *workers)
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	} else {
        memc["idfa"] = *idfaOpt
        memc["gaid"] = *gaidOpt
        memc["adid"] = *adidOpt
        memc["dvid"] = *dvidOpt
	}
}

func start(workers int, pattern string, memc map[string]string) {
	return
}
