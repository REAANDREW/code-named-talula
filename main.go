package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/satori/go.uuid"
)

var endpoints = map[uuid.UUID]Endpoint{}

func main() {
	var wait sync.WaitGroup
	wait.Add(1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Going to stop now...")
		wait.Done()
		os.Exit(1)
	}()

	StartV8Worker()

	go StartAdminServer()

	go StartTransformServer()

	wait.Wait()
}
