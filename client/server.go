package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./build/")))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("Finally, there's an error %s", err)
	}
}
