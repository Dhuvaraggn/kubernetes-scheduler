
package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const schedulerName = "hightower"

func main() {
	log.Println("Starting custom scheduler...")

	doneChan := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go monitorUnscheduledPods(doneChan, &wg)

	wg.Add(1)
	go reconcileUnscheduledPods(30, doneChan, &wg)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			log.Printf("Shutdown signal received, exiting...")
			close(doneChan)
			wg.Wait()
			os.Exit(0)
		}
	}
}
