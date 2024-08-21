package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

const serviceAddr = "/tmp/service.sock4"

func main() {
	ctx := context.Background()

	conn, err := net.Dial("unix", serviceAddr)
	if err != nil {
		log.Fatalf("Dial %q: %v", serviceAddr, err)
	}
	defer conn.Close()

	cli := jrpc2.NewClient(channel.Line(conn, conn), nil)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter messages to send to the server. Type 'exit' to quit.")

	for {
		// Prompt the user for input
		fmt.Print("Message: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Check for exit command
		if input == "exit" {
			fmt.Println("Exiting...")
			break
		}

		// Send the input to the server
		var amountOfChar int
		if err := cli.CallResult(ctx, "Count", []string{input}, &amountOfChar); err != nil {
			log.Printf("Call failed: %v", err)
			continue
		}

		// Print the result from the server
		fmt.Printf("Server response: %d characters\n", amountOfChar)
	}

	cli.Close()
}
