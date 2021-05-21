package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var addr = flag.String("addr", ":8081", "Address of Website")
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("/go/src/app/web/public"))))
	log.Println("Address of Website: ", addr)
	http.ListenAndServe(*addr, mux)
}
