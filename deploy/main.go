package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	binpath := os.Getenv("GOPATH")
	if len(binpath) == 0 {
		panic("GOPATH not set")
	}
	binpath = strings.Split(binpath, ";")[0]
	if len(binpath) == 0 {
		panic("problem parsing GOPATH")
	}
	binpath = filepath.Join(binpath, "bin")

	log.Printf("Serving from %s...", binpath)
	log.Panic(http.ListenAndServe(":3000", http.FileServer(http.Dir(binpath))))
}
