package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("heavy process starts")
	time.Sleep(10 * time.Second)
	log.Println("done")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("hello\n"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("starting server")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Received TERM signal.")

	gracefulShutdownTimeout, err := strconv.Atoi(os.Getenv("GRACEFUL_SHUTDOWN_TIMEOUT"))
	if err != nil {
		log.Printf("Convert graceful shutdown timeout failed: %v", err)
	}
	ctxShutdown, cancel := context.WithTimeout(context.Background(), time.Duration(gracefulShutdownTimeout)*time.Second)
	defer cancel()
	<-ctxShutdown.Done()
	log.Printf("Waiting up to %d s before terminating.", gracefulShutdownTimeout)
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server Shutdown:", err)
	} else {
		log.Println("Server was shutdown")
	}
}
