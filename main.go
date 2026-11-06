package main

import (
	"github.com/72sevenzy2/grouter/grouter"
	"net/http"
)

func main() {
	r := grouter.NewGrouter()

	r.Get("ping", func(w http.ResponseWriter, r *grouter.Request) {
		w.Write([]byte("pong"))
	})

	http.ListenAndServe(":9000", r)
}
