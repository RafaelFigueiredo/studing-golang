/*
https://medium.com/@thejasbabu/concurrency-in-go-e4a61ec96491
"It also shows how the select statement can be used to choose one out of several communications. A select statement chooses which of a set of possible send or receive operations will proceed. It looks similar to a switch statement but with the cases all referring to communication operations."
*/

package main

import (
	"fmt"
	"sync"
)

func main() {
	people := []string{"Anna", "Bob", "Jack", "Jill", "Dave", "Cody"}
	match := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(len(people))
	for _, name := range people {
		go Match(name, match, wg)
	}
	wg.Wait()
}

// Match either sends or receives, whichever possible, a name on the match
// channel and notifies the wait group when done.
func Match(name string, match chan string, wg *sync.WaitGroup) {
	select {
	case peer := <-match:
		fmt.Printf("%s sent a message to %s.\n", peer, name)
	case match <- name:
		// Wait for someone to receive my message.
	}
	wg.Done()
}
