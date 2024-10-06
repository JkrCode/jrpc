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

//unix file address relative to root
const serviceAddr = "/tmp/service.sock"

func main() {
	//context used to propagate cancellation input
	ctx := context.Background()

	// net.Dial is a low level function to establish a network connections (here not because
	// local) it returns a conn object that is used to read from and write to the other side
	// unix used to connect process on the same mashine
	// Unix sockets are more performant for inter process communication as it bypasses the network stack
	// no overhead when communicating
	conn, err := net.Dial("unix", serviceAddr)
	if err != nil {
		log.Fatalf("Dial %q: %v", serviceAddr, err)
	}
	defer conn.Close()

	//initialisieurng des jrpc2 Clients
	//channel.Line: framing to ensure server and client know when messages end etc. to ensure correct input/output format
	//conn conn say incoming and outgoing messages are written to the same connection
	//connection uses go channel under the hood
	//channel.line represents a configuration of the client 
	cli := jrpc2.NewClient(channel.Line(conn, conn), nil)

	//io reader for terminal input, that its coming from terminal is parameterized via os.Stin
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter messages to send to the server. Type 'exit' to quit.")

	for {
		// Prompt the user for input, written into the reader 
		fmt.Print("Message: ")
		input, _ := reader.ReadString('\n') //read until hit enter, ReadString is blocking
		input = strings.TrimSpace(input) // remove trailing whitespace

		// Check for exit command
		if input == "exit" {
			fmt.Println("Exiting...")
			break
		}

		// Send the input to the server
		// converts input to slice of string (because interface wants it that way)
		// with &amountOfChar binding the result to int var to capture response
		// ctx not directly interacted with, but we could have used ctx.withTimeout to 
		// manage the cancellation timing of the CallResult
		// Call Result is blocking 
		//Aufbau einer JSON-RPC Nachricht
		// {
  		// "jsonrpc": "2.0",
  		// "method": "subtract",
  		// "params": [42, 23],
  		// "id": 1
		// }
		//vorteile: mehrere Methodenaufrufe auf einmal übergeben, besserer fit wenn es darum geht 
		//methoden aufzurufen anstatt ressourcen zu manipulieren (wenn crud kein sinn macht)
		// bsp. substract ist eine method und keine ressource, also macht ein post /substract eventuell wenig intuitiv sinn
		var amountOfChar int
		if err := cli.CallResult(ctx, "Count", []string{input}, &amountOfChar); err != nil {
			log.Printf("Call failed: %v", err)
			continue //next loop
		}

		// Print the result from the server
		fmt.Printf("Server response: %d characters\n", amountOfChar)
	}

	cli.Close()
}
