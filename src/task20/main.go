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

	var ch = make(chan duration)
	go readDuration("src/data.txt", ch)

	var wg sync.WaitGroup
	for n := 0; ; n++ {
		var dur = <-ch
		if dur.h == 0 && dur.m == 0 && dur.s == 0 {
			break
		}

		wg.Add(1)
		go runTask(dur, &wg)
	}
	wg.Wait()
}

func readDuration(filePath string, ch chan duration) {
	var file, e = os.Open(filePath)
	if e != nil {
		panic(e)
	}
	defer file.Close()

	for {
		var dur duration
		var _, e = fmt.Fscanf(file, "%dh%dm%ds\n", &dur.h, &dur.m, &dur.s)

		if e == io.EOF {
			ch <- duration{0, 0, 0}
		}

		ch <- dur
	}
}

func runTask(dur duration, wg *sync.WaitGroup) {
	log.Printf("task started\n")
	time.Sleep(1000 * time.Millisecond)
	log.Printf("task finished (%dh%dm%ds)", dur.h, dur.m, dur.s)
	wg.Done()
}
