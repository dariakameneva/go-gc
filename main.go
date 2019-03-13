package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Below is an example of using our PrintMemUsage() function
	// Print our starting memory usage (should be around 0mb)
	PrintMemUsage()

	finish := make(chan struct{})
	go listenForSignals(finish)

	var overall [][]int
	defer func() {
		// Clear our memory and print usage, unless the GC has run 'Alloc' will remain the same
		overall = nil
		PrintMemUsage()

		// Force GC to clear up, should see a memory drop
		runtime.GC()
		PrintMemUsage()
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(done <-chan struct{}) {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
			}
			// Allocate memory using make() and append to overall (so it doesn't get
			// garbage collected). This is to create an ever increasing memory usage
			// which we can track. We're just using []int as an example.
			a := make([]int, 0, 999999)
			overall = append(overall, a)

			// Print our memory usage at each interval
			PrintMemUsage()
			time.Sleep(time.Second)
		}
	}(finish)

	wg.Wait()
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func listenForSignals(done chan<- struct{}) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for sig := range signalChan {
		switch sig {
		default:
			log.Panicf("Unexpected signal: %v", sig)

		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			close(done)
		}
	}
}
