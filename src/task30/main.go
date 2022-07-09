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

type task struct {
	dur duration
	num int
}

type queue struct {
	tasks chan task
	mu sync.Mutex
}

func main() {
	var logFile, _ = os.OpenFile("log/default.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666,
	)
	defer logFile.Close()
	log.SetOutput(logFile)

	var capacity uint64
	fmt.Printf("enter queue capacity: ")
	fmt.Fscanf(os.Stdin, "%d", &capacity)

	var q = queue {tasks: make(chan task, capacity)}

	go readDuration("src/data.txt", &q)
	runTaskQueue(&q)
}

func readDuration(filePath string, q *queue) {
	var file, e = os.Open(filePath)
	if e != nil {
		panic(e)
	}
	defer file.Close()

	for n := 1;; n++ {
		var dur duration
		var _, e = fmt.Fscanf(file, "%dh%dm%ds\n", &dur.h, &dur.m, &dur.s)

		for len(q.tasks) == cap(q.tasks) {}
		q.mu.Lock()

		if e == io.EOF {
			q.tasks <- task{}
			q.mu.Unlock()
			break
		}

		q.tasks <- task {dur, n}
		log.Printf("push task (%d/%d)\n", len(q.tasks), (cap(q.tasks)))
		q.mu.Unlock()
	}
}

func runTask(t task, wg *sync.WaitGroup) {
	var dur = &t.dur

	log.Printf("task #%d started\n", t.num)
	time.Sleep(100 * time.Millisecond)
	log.Printf("task #%d finished (%dh%dm%ds)", t.num, dur.h, dur.m, dur.s)

	wg.Done()
}

func runTaskQueue(q *queue) {
	var wg sync.WaitGroup

	for {
		for len(q.tasks) == 0 {}
		q.mu.Lock()

		var t = <-q.tasks
		var dur = &t.dur
		if dur.h == 0 && dur.m == 0 && dur.s == 0 {
			break
		}
		log.Printf("pop task (%d/%d)\n", len(q.tasks), cap(q.tasks))
		
		q.mu.Unlock()

		wg.Add(1)
		go runTask(t, &wg)
	}

	wg.Wait()
}
