package http_funcs

import "net/http"

func FileServer() {
	fileServer := http.FileServer(http.Dir("./sample-dir/"))
	mux := http.NewServeMux()
	mux.Handle("/", fileServer)
	http.ListenAndServe(":8990", mux)
}
