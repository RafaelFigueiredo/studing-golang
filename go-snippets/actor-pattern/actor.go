/*
	https://www.appdynamics.com/blog/engineering/three-productive-go-patterns-put-radar/
*/

package pattern

type Actor struct {
	eventc chan Event
	requestc chan reqRes
	quitc chan struct{}
}

func (a *Actor) loop() {
	for {
		select {
		case e := <-eventc:
			a.consumeEvent(e)
		case r := <-requestc:
			res, err := a.handleRequest(r.req)
			r.resc <- resErr{res, err}
		case <-quitc:
			return
		}
	}
}

// Finally, we push onto those channels in our exported methods, forming our public API, which is naturally goroutine-safe.

func (a *Actor) SendEvent(e Event) {
	a.eventc <- e
}

func (a *Actor) MakeRequest(r *Request) (*Response, error) {
	resc := make(chan resErr)
	a.requestc <- reqRes{req: r, resc: resc}
	res := <-resc
	return res.res, res.err
}

func (a *Actor) Stop() {
	close(a.quitc)
}

type reqRes struct {
	req *Request
	resc chan resErr

type resErr struct {
	res *Response
	err error
}
