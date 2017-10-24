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
	//"strings"
	//"time"
	"strings"
	"errors"
	//"github.com/golang/protobuf/proto"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/livestalker/python-17-hw11/appsinstalled"
	"strconv"
	"github.com/golang/protobuf/proto"
	"sync"
	"sync/atomic"
)

var NORMAL_ERR_RATE = 0.01

type Task struct {
	key   string
	value []byte
}

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
	var wg sync.WaitGroup
	files, err := filepath.Glob(*pattern)
	if err != nil {
		log.Fatal(err)
	}
	sort.Strings(files)
	for _, f := range files {
		wg.Add(1)
		go handle_file(f, memc, &wg)
	}
	wg.Wait()
}

func handle_file(filename string, memc *map[string]string, wg *sync.WaitGroup) {
	var _processed uint64 = 0
	var _errors uint64 = 0
	var doneFlag sync.WaitGroup
	defer wg.Done()
	mClients := make(map[string]*memcache.Client)
	taskCh := make(map[string](chan *Task))
	for key, value := range *memc {
		mClients[key] = memcache.New(value)
		taskCh[key] = make(chan *Task)
		doneFlag.Add(1)
		go insert_appsinstalled(mClients[key], taskCh[key], &_processed, &_errors, &doneFlag)
	}
	log.Printf("Processing: %s file.", filename)
	fh, err := os.Open(filename)
	if err != nil {
		log.Printf("File: %s, error: %s", filename, err)
		return
	}
	defer fh.Close()
	gz, err := gzip.NewReader(fh)
	if err != nil {
		log.Println(err)
		return
	}
	defer gz.Close()
	scanner := bufio.NewScanner(gz)

	for scanner.Scan() {
		line := scanner.Text()
		devType, devId, bytes, err := parse_appsinstalled(line)
		key := fmt.Sprintf("%s:%s", devType, devId)
		taskCh[devType] <- &Task{
			key:   key,
			value: bytes,
		}
		if err != nil {
			log.Printf("Line: %s, error: %s", line, err)
		}
	}
	for _, value := range taskCh {
		close(value)
	}
	doneFlag.Wait()
	totalProcessed := atomic.LoadUint64(&_processed)
	totalErrors := atomic.LoadUint64(&_errors)
	log.Printf("Total lines %d.", totalProcessed + totalErrors)
	if totalProcessed == 0 {
		log.Printf("File %s did not processsed", filename)
		return
	}
	errRate := float64(totalErrors) / float64(totalProcessed)
	if errRate < NORMAL_ERR_RATE {
		log.Printf("Acceptable error rate (%f). Successfull load.", errRate)
		//TODO Rename file
	} else {
		log.Printf("High error rate (%f > %f). Failed load.", errRate, NORMAL_ERR_RATE)
	}
}

func parse_appsinstalled(line string) (string, string, []byte, error) {
	var apps []uint32
	parts := strings.Split(strings.TrimSpace(line), "\t")
	if len(parts) != 5 {
		return "", "", nil, errors.New("error in format\n")
	}
	devType := parts[0]
	devId := parts[1]
	lat, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return "", "", nil, errors.New("float parsing error")
	}

	lon, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return "", "", nil, errors.New("float parsing error")
	}
	for _, el := range strings.Split(parts[4], ",") {
		app, err := strconv.ParseUint(el, 10, 32)
		if err != nil {
			continue
		}
		apps = append(apps, uint32(app))
	}
	ua := appsinstalled.UserApps{
		Lat:  &lat,
		Lon:  &lon,
		Apps: apps,
	}
	bytes, err := proto.Marshal(&ua)
	if err != nil {
		return "", "", nil, errors.New("marshaling error")
	}
	return devType, devId, bytes, nil
}

func insert_appsinstalled(mc *memcache.Client, tasks <-chan *Task, _processed *uint64, _errors *uint64, doneFlag *sync.WaitGroup) {
	defer doneFlag.Done()
	for t := range tasks {
		err := mc.Set(&memcache.Item{Key: t.key, Value: t.value})
		if err != nil {
			log.Println(err)
			atomic.AddUint64(_errors, 0)
		}
		atomic.AddUint64(_processed, 1)
	}
}
