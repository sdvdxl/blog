package main

import (
	"net/http"
	"log"
	"os/exec"
)

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	//用于git的webhook，触发pull
	http.HandleFunc("/_blog/_pull", func(writer http.ResponseWriter, request *http.Request) {
		cmd := exec.Command("git","pull")
		if err:=cmd.Start(); err!=nil {
			log.Println("git pull error", err)
		}
	})

	go func() {
		if err := http.ListenAndServeTLS(":443", "todu.crt", "todu.top.key", nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	err := http.ListenAndServe(":80", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}