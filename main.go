package main

import (
	"context"
	"fmt"
	"log"
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
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// graceful shutdown
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan
	fmt.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
