package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "you monitored my API...\n")
}

func main() {
	// create a server multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)

	// start the server
	// nil - use the default server multiplexer
	err := http.ListenAndServe("localhost:8000", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
