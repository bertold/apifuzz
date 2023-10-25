package main

import (
	"bufio"
	"context"
	"net/http"
	"os"
	"time"
)

func main() {
	s := StartServer()
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	s.Shutdown(context.Background())
}

func StartServer() *http.Server {
	http.HandleFunc("/date", DateHandler)
	http.HandleFunc("/reverse", ReverseHandler)
	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func(server *http.Server) {
		server.ListenAndServe()
	}(s)
	return s
}
