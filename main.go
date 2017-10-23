package main

import (
	"log"
	"net/http"

	"shortURL/core"
)

func main() {

	if err := core.LoadConfg(); err != nil {
		panic(err)
	}

	core.LinkDB()

	http.HandleFunc("/", core.GetOriginalURL)
	http.HandleFunc("/_genShortURL", core.GenShortURL)

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
