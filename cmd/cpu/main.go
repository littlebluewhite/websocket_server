package main

import (
	"fmt"
	"sync"
)

var counter = 0
var wg = sync.WaitGroup{}

func AddCounter() {
	defer wg.Done()
	counter++
}

func main() {
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		go AddCounter()
	}

	wg.Wait()
	fmt.Println(counter) // ?
}
