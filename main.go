package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/borowiak-m/interview-assignment-week-17/handlers"
)

func main() {
	// mongodb+srv://challengeUser:WUMgIwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true

	// handlers
	sm := http.NewServeMux()
	// 	Records in memory handler
	inmemHandler := handlers.NewRecsInMem()

	// register routes
	// /GET records from Mongo
	// /GET from map / records in memory
	sm.Handle("/inmemory", inmemHandler)
	// /POST to map / records in memory

	// define server
	s := &http.Server{
		Addr:         ":3000",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	// non-blocking listen and serve
	go func() {
		fmt.Println("Starting server on port 3000")

		err := s.ListenAndServe()
		if err != nil {
			fmt.Println("Error starting server: $s\n", err)
			os.Exit(1)
		}
	}()

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan
	fmt.Println("Received terminate, graceful shutdown", sig)
	// define timeout context (cancel not used here, hence _ )
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
