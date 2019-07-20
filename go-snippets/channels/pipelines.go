package main
// https://blog.golang.org/pipelines


/* Consider a pipeline with three stages.
The first stage, gen, is a function that converts a list of integers to a channel that emits the integers in the list. The gen function starts a goroutine that sends the integers on the channel and closes the channel when all the values have been sent: */
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}


/* The second stage, sq, receives integers from a channel and returns a channel that emits the square of each received integer. After the inbound channel is closed and this stage has sent all the values downstream, it closes the outbound channel: */
func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

/* The main function sets up the pipeline and runs the final stage: it receives values from the second stage and prints each one, until the channel is closed:
 */
func main() {
    // Set up the pipeline.
    c := gen(2, 3)
    out := sq(c)

    // Consume the output.
    fmt.Println(<-out) // 4
    fmt.Println(<-out) // 9
}


/* Since sq has the same type for its inbound and outbound channels, we can compose it any number of times. We can also rewrite main as a range loop, like the other stages: */

func main() {
    // Set up the pipeline and consume the output.
    for n := range sq(sq(gen(2, 3))) {
        fmt.Println(n) // 16 then 81
    }
}