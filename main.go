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
	"strings"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/livestalker/python-17-hw11/appsinstalled"
	"strconv"
	"github.com/golang/protobuf/proto"
	"sync"
)

var NORMAL_ERR_RATE = 0.01

type Task struct {
	key   string
	value []byte
}

type Result struct {
	processed int
	errors    int
}

var memc map[string] string
var pattern *string
var workers *int

func init() {
	memc = make(map[string]string)
	pattern = flag.String("pattern", "", "Files pattern.")
	memc["idfa"] = *flag.String("idfa", "127.0.0.1:33013", "Memcache for idfa.")
	memc["gaid"] = *flag.String("gaid", "127.0.0.1:33014", "Memcache for gaid.")
	memc["adid"] = *flag.String("adid", "127.0.0.1:33015", "Memcache for adid.")
	memc["dvid"] = *flag.String("dvid", "127.0.0.1:33016", "Memcache for dvid.")
	workers = flag.Int("workers", runtime.NumCPU(), "Count of forkers.")
}

func main() {
	var wg sync.WaitGroup
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	}
	if *pattern == "" {
		log.Fatal("Pattern not found in arguments")
	}
	runtime.GOMAXPROCS(*workers)
	files, err := filepath.Glob(*pattern)
	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
	sort.Strings(files)
	for _, f := range files {
		wg.Add(1)
		go handleFile(f, memc, &wg)
	}
	wg.Wait()
	for _, f := range files {
		renameFile(f)
		if err != nil {
			log.Print(err)
		}
		log.Printf("Rename %s file.", f)
	}
}

func handleFile(filename string, memc map[string]string, wg *sync.WaitGroup) {
	var results = make(chan * Result, len(memc))
	var doneFlag sync.WaitGroup
	var parserFlag sync.WaitGroup
	var clients int
	defer wg.Done()
	mClients := make(map[string]*memcache.Client)
	taskCh := make(map[string](chan *Task))
	for key, value := range memc {
		mClients[key] = memcache.New(value)
		taskCh[key] = make(chan *Task)
		doneFlag.Add(1)
		go insertAppsinstalled(mClients[key], taskCh[key], results, &doneFlag)
		clients++
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
		parserFlag.Add(1)
		go parseAppsinstalled(line, taskCh, &parserFlag)
	}
	parserFlag.Wait()
	for _, value := range taskCh {
		close(value)
	}
	var totalProcessed int
	var totalErrors int
	doneFlag.Wait()
	close(results)
	for r := range results {
		totalProcessed += r.processed
		totalErrors += r.errors
	}
	log.Printf("Total lines %d if file %s.", totalProcessed+totalErrors, filename)
	if totalProcessed == 0 {
		log.Printf("File %s did not processsed.", filename)
		return
	}
	errRate := float64(totalErrors) / float64(totalProcessed)
	if errRate < NORMAL_ERR_RATE {
		log.Printf("Acceptable error rate (%f). Successfull load file %s.", errRate, filename)
	} else {
		log.Printf("High error rate (%f > %f). Failed load file %s.", errRate, NORMAL_ERR_RATE, filename)
	}
}

func parseAppsinstalled(line string, taskCh map[string](chan *Task), parserFlag *sync.WaitGroup) {
	var apps []uint32
	defer parserFlag.Done()
	parts := strings.Split(strings.TrimSpace(line), "\t")
	if len(parts) != 5 {
		log.Printf("Line: %s, error: %s", line, "error in format.")
	}
	devType := parts[0]
	devId := parts[1]
	lat, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		log.Printf("Line: %s, error: %s", line, "float parsing error.")
	}

	lon, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		log.Printf("Line: %s, error: %s", line, "float parsing error.")
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
		log.Printf("Line: %s, error: %s", line, "marshaling error")
	}
	key := fmt.Sprintf("%s:%s", devType, devId)
	taskCh[devType] <- &Task{
		key:   key,
		value: bytes,
	}
}

func insertAppsinstalled(mc *memcache.Client, tasks <-chan *Task, results chan<- *Result, doneFlag *sync.WaitGroup) {
	var _errors, _processed int
	defer doneFlag.Done()
	for t := range tasks {
		err := mc.Set(&memcache.Item{Key: t.key, Value: t.value})
		if err != nil {
			log.Println(err)
			_errors ++
		} else {
			_processed ++
		}
	}
	results <- &Result{
		processed: _processed,
		errors: _errors,
	}
}

func renameFile(filename string) error {
	newFilename := filepath.Dir(filename) + "/." + filepath.Base(filename)
	return os.Rename(filename, newFilename)
}
