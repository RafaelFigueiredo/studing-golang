package main

import (
	"fmt"
	"sync"
	"time"
)

func gen(done <-chan struct{}, numbers ...int) <-chan int {
	out := make(chan int, len(numbers))

	for _, n := range numbers {
		out <- n
	}

	close(out)
	return out

	/*go func() {
		for _, n := range numbers {
			out <- n
		}
		close(out)
	}()*/

}

func sg(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}

		}
	}()

	return out
}

func main() {
	t_begin := time.Now()

	// Create a done channel to broadcast end of process
	done := make(chan struct{})
	defer close(done)

	// Feed the queue with information to be processed
	in := gen(done, 2, 3)

	// Distribute the sq work acros two go routines that both read from in
	c1 := sg(done, in)
	c2 := sg(done, in)

	// Consumed the merge output from c1 and c2
	for res := range merge(done, c1, c2) {
		fmt.Println(res)
	}

	fmt.Printf("Time Elapsed:: %d ns", time.Since(t_begin).Nanoseconds())
}

func merge(done <-chan struct{}, ch ...<-chan int) <-chan int {
	// Create our output channel and declare the WaitGroupt that will be used to ensure that all goroutines
	// run befere close the output channel
	var wg sync.WaitGroup
	out := make(chan int, 1)

	// Our function that will merge inbound information from all channels to one
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}

	}

	// Set WaitGroup
	wg.Add((len(ch)))

	// Start our merge goroutine
	for _, c := range ch {
		go output(c)
	}

	// Ensure that all go routines run before close the channel
	go func() {
		wg.Wait()
		close(out)
	}()

	// Return the output channel
	return out
}
