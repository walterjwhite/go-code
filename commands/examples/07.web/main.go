package main

import (
	"log"
	
	"net/http"
	"github.com/ddo/rq"
)

func main() {
	r := rq.Get("https://pnc.com")

	// send with golang default HTTP client
	req, err := r.ParseRequest()
	
	if err != nil {
		panic(err)
	}
	
	res, err := http.DefaultClient.Do(req)
	
	if err != nil {
		panic(err)
	}
	
	defer res.Body.Close()
	
	log.Printf("Response:\n%v\n", res)
}
