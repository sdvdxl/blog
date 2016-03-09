package main

import (
	"net/http"
	"log"
)

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}