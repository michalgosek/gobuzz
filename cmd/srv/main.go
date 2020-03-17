package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gobuzz/pkg/domain/adding"
	"github.com/gobuzz/pkg/domain/responding"
	"github.com/gobuzz/pkg/http/rest"
	"github.com/gobuzz/pkg/storage/memory"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Initializing storage and services.
	s := new(memory.ResponseFetch)
	adder := adding.NewService(&s.Fetches)        // adding service
	respsr := responding.NewService(&s.Responses) // responsing service (for Gopher)

	srv := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           rest.ServHandler(adder, respsr),
		MaxHeaderBytes:    1 << 20, //1MB
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Println("GoBuzz server is running: http://localhost:8080")
	return srv.ListenAndServe()
}
