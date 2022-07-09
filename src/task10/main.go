package main

import (
	"fmt"
	"io"
	"log"
	"os"
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

	var durs = readDuration("src/data.txt")
	for n, dur := range durs {
		log.Printf("task %d started\n", n)
		runTask(dur)
		log.Printf("task %d finished (%dh%dm%ds)", n, dur.h, dur.m, dur.s)
	}
}

func readDuration(filePath string) []duration {
	var durs []duration
	var file, err = os.Open(filePath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	for {
		var dur duration
		var _, e = fmt.Fscanf(file, "%dh%dm%ds\n", &dur.h, &dur.m, &dur.s)

		if e == io.EOF {
			break
		}

		durs = append(durs, dur)
	}

	return durs
}

func runTask(dur duration) {
	time.Sleep(1000 * time.Millisecond)
}
