/* "It will take in a request and provide you its HTTP 1.1 request format. It will dump out all the headers,
    parameters and even the body (if you provide a true value for the body parameter)."
	https://rominirani.com/golang-tip-capturing-http-client-requests-incoming-and-outgoing-ef7fcdf87113?gi=9a38214099eb*/
	
package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httputil"
)

func DumpRequest(w http.ResponseWriter, req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Fprint(w, string(requestDump))
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/dumprequestG", DumpRequest).Methods("GET")
	router.HandleFunc("/dumprequestP", DumpRequest).Methods("POST")
	log.Fatal(http.ListenAndServe(":12345", router))
}