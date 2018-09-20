package main

import (
	"fmt"
	"strconv"
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

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for i := range in {
		wg.Add(1)
		go singleHashJob(i, out, wg, mu)
	}
}

func singleHashJob(in interface{}, out chan interface{}, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	data := strconv.Itoa(in.(int))

	mu.Lock()
	md5Data := DataSignerMd5(data)
	mu.Unlock()

	crc32Chan := make(chan string)
	go asyncCrc32Signer(data, crc32Chan)

	crc32Data := <-crc32Chan
	crc32Md5Data := DataSignerCrc32(md5Data)

	out <- crc32Data + "~" + crc32Md5Data
}

func asyncCrc32Signer(data string, out chan string) {
	out <- DataSignerCrc32(data)
}

func main() {
	fmt.Println("pew")
}

// сюда писать код
