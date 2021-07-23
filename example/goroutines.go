package main

import (
	"fmt"
	"strconv"
	"sync"
)

func worker(wg *sync.WaitGroup, cs chan string, i int) {
	defer wg.Done()
	cs <- "worker" + strconv.Itoa(i)
}

func monitorWorker(wg *sync.WaitGroup, cs chan string) {
	wg.Wait()
	close(cs)
}
func main() {
	wg := &sync.WaitGroup{}
	cs := make(chan string)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(wg, cs, i)
	}

	go monitorWorker(wg, cs)

	for i := range cs {
		fmt.Println(i)

	}
}
