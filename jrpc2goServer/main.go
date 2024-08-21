package main

import (
	"context"
	"fmt"
	"jrpcServer/methods"
	"log"
	"net"
	"sync"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
)

// CountService provides a method to count characters in a string.
type CountService struct {
	mu         sync.Mutex
	totalCount int
}

const serviceAddr = "/tmp/service.sock3"

func main() {
	// Start the server listening on the local network.
	lst, err := net.Listen(jrpc2.Network(serviceAddr))
	if err != nil {
		log.Fatalf("Listen %q: %v", serviceAddr, err)
	}
	defer lst.Close()

	// Set up a service with some trivial methods
	svc := server.Static(handler.Map{
		"Join":  handler.New(methods.Join),
		"Log":   handler.NewPos(methods.LogPrint, "level", "message"),
		"Count": handler.New(methods.Count),
	})
	ctx := context.Background()
	go server.Loop(ctx, server.NetAccepter(lst, channel.Line), svc, nil)

	//Dial the server and set up a client.
	conn, err := net.Dial(jrpc2.Network(serviceAddr))
	if err != nil {
		log.Fatalf("Dial %q: %v", serviceAddr, err)
	}
	defer conn.Close()

	cli := jrpc2.NewClient(channel.Line(conn, conn), nil)

	// -------------------------------------------------------------------------
	// Post a notification.
	if err := cli.Notify(ctx, "Log", handler.Obj{
		"level":   "debug",
		"message": "I am in your forest, logging your logs",
	}); err != nil {
		log.Fatalf("Notify failed: %v", err)
	}

	// -------------------------------------------------------------------------
	// Issue a call with a response.
	var result string
	if err := cli.CallResult(ctx, "Join", []string{"hello", "beautiful", "world"}, &result); err != nil {
		log.Fatalf("Call failed: %v", err)
	}
	fmt.Println(result)

	var amountOfChar int
	if err := cli.CallResult(ctx, "Count", []string{"hello world"}, &amountOfChar); err != nil {
		log.Fatalf("Call failed: %v", err)
	}
	fmt.Println(amountOfChar)

	cli.Close()

}
