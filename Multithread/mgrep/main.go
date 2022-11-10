package main

import (
	"fmt"
	"mgrep/worker"
	"mgrep/worklist"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
)

func discoverDirs(wl *worklist.WorkList, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Readdir error:", err)
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			nextPath := filepath.Join(path, entry.Name())
			discoverDirs(wl, nextPath)
		} else {
			wl.Add(worklist.NewJob(filepath.Join(path, entry.Name())))
		}
	}
}

var args struct {
	SearchTerm string `arg:"positional,required"`
	SearchDir  string `arg:"positional"`
}

func main() {
	arg.MustParse(&args)
	var workerWg sync.WaitGroup

	wl := worklist.New(100)
	results := make(chan worker.Result, 100)
	numWorkers := 10
	workerWg.Add(1)
	go func() {
		defer workerWg.Done()
		discoverDirs(&wl, args.SearchDir)
		wl.FinalSize(numWorkers)
	}()
	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for {
				workEntry := wl.Next()
				if workEntry.Path != "" {
					wokerResult := worker.FindInFile(workEntry.Path, args.SearchTerm)
					if wokerResult != nil {
						for _, v := range wokerResult.Inner {
							results <- v
						}
					}
				} else {
					return
				}
			}
		}()
	}
	blockWorkersWg := make(chan struct{})
	go func() {
		workerWg.Wait()
		close(blockWorkersWg)
	}()
	var displayWg sync.WaitGroup
	displayWg.Add(1)
	go func() {
		for {
			select {
			case r := <-results:
				fmt.Printf("%v[%v]:%v\n", r.Path, r.LineNum, r.Line)
			case <-blockWorkersWg:
				if len(results) == 0 {
					displayWg.Done()
					return
				}
			}
		}
	}()
	displayWg.Wait()
}
