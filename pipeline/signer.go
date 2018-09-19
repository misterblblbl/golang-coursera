package main

import (
	"fmt"
	"sync"
)

var concurrency = 6

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})

	for _, job := range jobs {
		wg.Add(1)

		out := make(chan interface{})
		go jobWorker(job, in, out, wg)
		in = out
	}

	wg.Wait()
}

func jobWorker(job job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)

	job(in, out)
}

func SingleHash() {}

func main() {
	fmt.Println("pew")
}

// сюда писать код
