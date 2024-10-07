package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"jrpcServer/handlers"
	"jrpcServer/shared"

	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
)

const serviceAddr = "/tmp/service.sock"

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

	messageCounter := make(chan struct{})
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter Filter for server:")
	fmt.Print("Event Type:")
	filterEventType, _ := reader.ReadString('\n')
	filterEventType = strings.TrimSpace(filterEventType)

	fmt.Print("from dd/mm/yyyy to now: \n")
	filterTimeStamp := time.Now()

	fmt.Print("SourceIP:")
	filterSourceIP, _ := reader.ReadString('\n')

	fmt.Print("UserID (email):")
	filterUserId, _ := reader.ReadString('\n')

	fmt.Print("Severity Level:")
	filterSeverityString, _ := reader.ReadString('\n')
	filterSeverity, _ := strconv.Atoi(filterSeverityString)

	filter := shared.SecurityLog{
		TimeStamp: &filterTimeStamp,
		EventType: &filterEventType,
		SourceIP:  &filterSourceIP,
		UserID:    &filterUserId,
		Severity:  &filterSeverity,
	}

	svc := server.Static(handler.Map{
		"Count":  handler.New(handlers.Count(messageCounter)),
		"Filter": handler.New(handlers.Filter(filter, messageCounter)),
	})

	// Handle shutdown signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx := context.Background()
	go func() {
		fmt.Println("Server is running...")
		server.Loop(ctx, server.NetAccepter(lst, channel.Line), svc, nil)
	}()

	// Aggregate the total number of messages received
	totalMessages := 0
	go collectCounts(messageCounter, &totalMessages)

	<-c
	fmt.Println("\nServer is shutting down...")
	fmt.Printf("Total messages received: %d\n", totalMessages)
}

func Atoi(filterSeverityString string) {
	panic("unimplemented")
}

func collectCounts(messageCounter chan struct{}, totalMessages *int) {
	// ranging over a channel does mean to continiously pulling values from the channel and
	for range messageCounter {
		*totalMessages++
	}
}
