//go:build serve

package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("wasm"))

	log.Println("Serving wasm at http://localhost:8080")
	err := http.ListenAndServe(":8080", fs)
	if err != nil {
		log.Fatal(err)
	}
}
