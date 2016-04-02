package main

import "github.com/pkg4go/httprange"
import "github.com/justinas/alice"
import "net/http"
import "strconv"

func main() {
	var chain = alice.New(httprange.New()).Then(http.HandlerFunc(resText))

	http.ListenAndServe(":3000", chain)
}

func resText(res http.ResponseWriter, req *http.Request) {
	body := []byte("Hello, world!")
	size := len(body)

	res.Header().Set("Content-Length", strconv.Itoa(size))
	res.Write(body)
}
