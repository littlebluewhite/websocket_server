package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int, 2)
	go func() {
		c <- 1
		c <- 2
	}()
	fmt.Println(<-c)
	time.Sleep(1 * time.Second)
	close(c)
	fmt.Println(<-c)
	time.Sleep(3 * time.Second)
	fmt.Println("end")
}
