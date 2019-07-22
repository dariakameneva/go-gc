package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/trace"
	"strings"
	"sync"
)

const size = 1000

type tracker struct {
	sync.Mutex
	m map[string]struct{}
}

var t tracker

func main() {
	runtime.GOMAXPROCS(4)

	trace.Start(os.Stdout)
	defer trace.Stop()

	dir := "/Users/dariakameneva/gopath/src/"
	t.m = make(map[string]struct{})

	if err := filepath.Walk(dir, walk); err != nil {
		log.Fatal(err)
	}
}

func walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	go func() {
		file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Println(err)
			return
		}

		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			return
		}

		for _, l := range strings.Split(string(content), "\n") {
			if strings.Contains(l, "go func") {
				t.Lock()
				t.m[filepath.Base(path)] = struct{}{}
				t.Unlock()
			}
		}
	}()
	return nil
}
