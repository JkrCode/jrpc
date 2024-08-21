package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
)

// CountService provides a method to count characters in a string.
type CountService struct {
	mu         sync.Mutex
	totalCount int
}

const serviceAddr = "/tmp/service.sock4"

func main() {
	// Remove any existing socket file
	if _, err := os.Stat(serviceAddr); err == nil {
		os.Remove(serviceAddr)
	}

	// Start the server listening on the local network.
	lst, err := net.Listen("unix", serviceAddr)
	if err != nil {
		fmt.Printf("Listen %q: %v", serviceAddr, err)
		return
	}
	defer lst.Close()

	// Set up a service procedure
	svc := server.Static(handler.Map{
		"Count": handler.New(func(ctx context.Context, req []string) (int, error) {
			fmt.Println("Received request:", req)
			return len(req[0]), nil
		}),
	})

	ctx := context.Background()
	fmt.Println("Server is running...")
	server.Loop(ctx, server.NetAccepter(lst, channel.Line), svc, nil)
}
