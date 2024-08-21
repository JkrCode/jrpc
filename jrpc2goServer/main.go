package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
)

type CountService struct {
	counters [numShards]Counter
}

type Counter struct {
	value int64
	_     [120]byte // Padding (optimized for Mac ARM64)
}

const serviceAddr = "/tmp/service.sock"
const numShards = 64

func main() {
	// Remove any existing socket
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

	// Initialize the service with padded counters (initialized with 0 values, when nothing added)
	service := &CountService{}

	//procedure
	svc := server.Static(handler.Map{
		"Count": handler.New(func(ctx context.Context, req []string) (int, error) {
			shard := len(req[0]) % numShards
			counter := &service.counters[shard]

			atomic.AddInt64(&counter.value, 1)

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
	totalMessages := int64(0)
	for i := 0; i < numShards; i++ {
		totalMessages += service.counters[i].value
	}

	fmt.Printf("Total messages received: %d\n", totalMessages)
}
