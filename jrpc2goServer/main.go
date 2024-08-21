package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
)

// CountService provides a method to count characters in a string.
type CountService struct {
	counters []int
	mu       []sync.Mutex
	shards   int
}

const serviceAddr = "/tmp/service.sock4"
const numShards = 16

func main() {
	// Remove any existing socket file
	if _, err := os.Stat(serviceAddr); err == nil {
		os.Remove(serviceAddr)
	}

	// Start the server listening on the local network.
	lst, err := net.Listen("unix", serviceAddr)
	if err != nil {
		fmt.Printf("Listen %q: %v\n", serviceAddr, err)
		return
	}
	defer lst.Close()

	// Initialize the service with sharded counters
	service := &CountService{
		counters: make([]int, numShards),
		mu:       make([]sync.Mutex, numShards),
		shards:   numShards,
	}

	// Set up a service procedure
	svc := server.Static(handler.Map{
		"Count": handler.New(func(ctx context.Context, req []string) (int, error) {
			// Determine which shard to use
			shard := len(req[0]) % service.shards
			service.mu[shard].Lock()
			service.counters[shard]++
			service.mu[shard].Unlock()

			fmt.Println("Received request:", req)
			return len(req[0]), nil
		}),
	})

	// Handle interrupt signal (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Run the server in a separate goroutine
	ctx := context.Background()
	go func() {
		fmt.Println("Server is running...")
		server.Loop(ctx, server.NetAccepter(lst, channel.Line), svc, nil)
	}()

	// Wait for the interrupt signal
	<-c
	fmt.Println("\nServer is shutting down...")

	// Aggregate the total number of messages received
	totalMessages := 0
	for i := 0; i < service.shards; i++ {
		totalMessages += service.counters[i]
	}

	fmt.Printf("Total messages received: %d\n", totalMessages)
}
