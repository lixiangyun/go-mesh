package main

import (
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/proxy/cfg", ProxyCfgHandler)
	mux.HandleFunc("/server/query", ServerQueryHandler)
	mux.HandleFunc("/server/register", ServerRegisterHandler)

	http.ListenAndServe(":3001", mux)
}
