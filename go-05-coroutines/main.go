package main

import (
    "fmt"
    "sync"
    "time"
)

// SafeResults stores results from multiple goroutines safely
type SafeResults struct {
    mu      sync.Mutex
    results []string
}

func (r *SafeResults) Add(result string) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.results = append(r.results, result)
}

func (r *SafeResults) All() []string {
    r.mu.Lock()
    defer r.mu.Unlock()
    return r.results
}

// download simulates downloading a file
func download(id int, filename string, fileCh chan<- string, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("⬇  Downloading %s...\n", filename)
    time.Sleep(time.Duration(id*200) * time.Millisecond)   // simulate network delay
    fileCh <- filename   // send downloaded file to processor
    fmt.Printf("✅ Downloaded %s\n", filename)
}

// process simulates processing downloaded files
func process(fileCh <-chan string, done <-chan bool, results *SafeResults, wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        select {
        case file := <-fileCh:
            processed := fmt.Sprintf("processed_%s", file)
            results.Add(processed)
            fmt.Printf("⚙  Processed: %s\n", processed)

        case <-done:
            fmt.Println("🛑 Processor received shutdown signal")
            return
        }
    }
}

func main() {
    files := []string{"report.pdf", "image.png", "data.csv", "notes.txt", "config.json"}

    fileCh  := make(chan string, len(files))   // buffered channel for downloaded files
    done    := make(chan bool)                  // signal channel for shutdown
    results := &SafeResults{}

    var downloadWg sync.WaitGroup
    var processWg  sync.WaitGroup

    // Launch processor goroutine
    processWg.Add(1)
    go process(fileCh, done, results, &processWg)

    // Launch downloader goroutines concurrently
    fmt.Println("=== Starting Downloads ===")
    for i, file := range files {
        downloadWg.Add(1)
        go download(i+1, file, fileCh, &downloadWg)
    }

    // Wait for all downloads to complete
    downloadWg.Wait()
    fmt.Println("\n=== All Downloads Complete ===")

    // Give processor time to handle remaining files
    time.Sleep(500 * time.Millisecond)

    // Signal processor to shut down
    done <- true
    processWg.Wait()

    // Print all results
    fmt.Println("\n=== Final Results ===")
    for _, r := range results.All() {
        fmt.Println(" ", r)
    }
}
```

---

## Mental Model Summary
```
GOROUTINE         go fn()
                  → lightweight concurrent task
                  → killed when main() exits
                  → use WaitGroup to wait for them

CHANNEL           ch := make(chan Type)
                  ch <- value    → send (blocks until received)
                  <-ch           → receive (blocks until sent)
                  make(chan T,3) → buffered (holds 3 before blocking)
                  close(ch)      → signal no more values coming
                  for v := range ch → receive until closed

SELECT            select { case <-ch1: ... case <-ch2: ... default: ... }
                  → picks whichever channel is ready
                  → default runs immediately if nothing ready
                  → time.After() for timeouts
                  → done channel for cancellation

WAITGROUP         wg.Add(1) before launching goroutine
                  wg.Done() inside goroutine when finished (use defer!)
                  wg.Wait() to block until all are done

MUTEX             mu.Lock() before touching shared data
                  mu.Unlock() after (use defer!)
                  embed in struct for cleaner code
```

---

## How They All Work Together
```
Goroutines   →  do the concurrent work
Channels     →  pass data safely between goroutines
Select       →  handle multiple channels / timeouts / cancellation
WaitGroup    →  know when all goroutines are finished
Mutex        →  protect shared data from race conditions