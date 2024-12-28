package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "you monitored my API...\n")
}

func getRandom(w http.ResponseWriter, r *http.Request) {
	randomNum := rand.Int()
	io.WriteString(w, strconv.Itoa(randomNum))
}

func getNum(w http.ResponseWriter, r *http.Request) {
	num := 2
	io.WriteString(w, strconv.Itoa(num))
}

func getSurprise(w http.ResponseWriter, r *http.Request) {
	randChoice := rand.Int()
	if randChoice&1 > 0 {
		http.NotFound(w, r)
	} else {
		io.WriteString(w, "im running")
	}
}

func main() {
	// create a server multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/get-random", getRandom)
	mux.HandleFunc("/get-num", getNum)
	mux.HandleFunc("/get-surprise", getSurprise)

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
