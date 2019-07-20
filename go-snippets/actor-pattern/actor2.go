package pattern

type Actor struct {
	actionc chan func()
	quitc   chan struct{}
}

func (a *Actor) loop() {
	for {
		select {
		case f := <-actionc:
			f()
		case <-quitc:
			return
		}
	}
}

func (a *Actor) SendEvent(e Event) {
	a.actionc <- func() {
		a.consumeEvent(e)
	}
}

func (a *Actor) HandleRequest(r *Request) (res *Response, err error) {
	done := make(chan struct{})
	a.actionc <- func() {
		defer close(done) // outer func shouldn't return before values are set
		res, err = a.handleRequest(r)
	}
	<-done

}
