package main

/*
https://medium.com/@thejasbabu/concurrency-in-go-e4a61ec96491

Closing the channel is important before ranging over it, else will lead to panics. Once a channel has been closed, you cannot send a value on the channel, but you can still receive from the closed channel.
*/

import (
	"fmt"
	"strings"
	"sync"
)

func main() {
	names := []string{"Dave", "Steve", "Corry"}
	c := toUpperCase(names)
	for upper := range c {
		fmt.Println(upper)
	}
}

func toUpperCase(names []string) <-chan string {
	var wg sync.WaitGroup
	ch := make(chan string, len(names))
	wg.Add(len(names))
	for _, name := range names {
		go func(name string) {
			defer wg.Done()
			ch <- strings.ToUpper(name)
		}(name)
	}
	wg.Wait()
	close(ch)
	return ch
}
