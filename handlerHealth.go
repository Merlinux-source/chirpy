package main

import "net/http"

func handlerHealth(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")

	rw.WriteHeader(200)
	rw.Write([]byte("OK"))
}
