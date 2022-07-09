package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type duration struct {
	h uint8
	m uint8
	s uint8
}

func main() {
	var logFile, _ = os.OpenFile("log/default.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666,
	)
	defer logFile.Close()
	log.SetOutput(logFile)

	var capacity uint64
	fmt.Fscanf(os.Stdin, "%d", &capacity)
	var queue = make(chan duration, capacity)

	// var ch = make(chan duration)

	var wg sync.WaitGroup
	go readDuration("src/data.txt", queue)
	// go runTaskQueue(queue, &wg)

	// for n := 0; ; n++ {
	// 	queue <- ch
	// }

	wg.Wait()
}

func readDuration(filePath string, ch chan duration) {
	var file, e = os.Open(filePath)
	defer file.Close()

	for {
		var dur duration
		fmt.Fscanf(file, "%dh%dm%ds\n", &dur.h, &dur.m, &dur.s)

		if e == io.EOF {
			ch <- duration{}
		}

		ch <- dur
		log.Printf("task pushed into the queue (%d/%d)\n", len(ch), cap(ch))
	}
}

func runTask(dur duration, wg *sync.WaitGroup) {
	log.Printf("task started\n")
	time.Sleep(1000 * time.Millisecond)
	log.Printf("task finished (%dh%dm%ds)", dur.h, dur.m, dur.s)
	wg.Done()
}

func runTaskQueue(queue chan duration, wg *sync.WaitGroup) {
	for {
		var dur = <-queue
		if dur.h == 0 && dur.m == 0 && dur.s == 0 {
			break
		}

		wg.Add(1)
		go runTask(dur, wg)
	}
}
