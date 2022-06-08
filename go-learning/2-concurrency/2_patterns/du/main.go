package main

// The du variant includes cancellation:
// it terminates quickly when the user hits return.

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

//!+1
var gDone = make(chan struct{})

func cancelled() bool {
	select {
	case <-gDone:
		return true
	default:
		return false
	}
}

//!-1

func main() {
	defer func() {
		time.Sleep(5 * time.Second)
		fmt.Println("the number of goroutines: ", runtime.NumGoroutine())
	}()

	// Determine the initial directories.
	roots := os.Args[1:]
	if len(roots) == 0 {
		roots = []string{"."}
	}

	//!+2
	// Cancel traversal when input is detected.
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		close(gDone)
	}()
	//!-2

	// Traverse each root of the file tree in parallel.
	fileSizes := make(chan int64)
	var wg sync.WaitGroup
	for _, root := range roots {
		wg.Add(1)
		go walkDir(root, &wg, fileSizes)
	}
	go func() {
		wg.Wait()
		close(fileSizes)
	}()

	// Print the results periodically.
	tick := time.Tick(500 * time.Millisecond)
	var nfiles, nbytes int64
loop:
	//!+3
	for {
		select {
		case <-gDone:
			// Drain fileSizes to allow existing goroutines to finish.
			// 确保 walkDir 中的goroutine可以正常退出。因为walkDir只在开头处由检测是否取消，
			// 当取消操作发生时，部分goroutine已经在运行了
			for range fileSizes {
			}
			return
		case size, ok := <-fileSizes:
			// ...
			//!-3
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes) // final totals
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
//!+4
func walkDir(dir string, wg *sync.WaitGroup, fileSizes chan<- int64) {
	defer wg.Done()
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			wg.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, wg, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

//!-4

var gSema = make(chan struct{}, 20) // concurrency-limiting counting semaphore

// dirents returns the entries of directory dir.
//!+5
func dirents(dir string) []os.FileInfo {
	select {
	case gSema <- struct{}{}: // acquire token
	case <-gDone:
		return nil // cancelled
	}
	defer func() { <-gSema }() // release token

	// ...read directory...
	//!-5

	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	defer f.Close()

	entries, err := f.Readdir(0) // 0 => no limit; read all entries
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		// Don't return: Readdir may return partial results.
	}
	return entries
}
