package main

import (
	"bufio"
	"context"
	"fmt"
	"jrpcClient/shared"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

// unix file address relative to root
const serviceAddr = "/tmp/service.sock"

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
		//Eventtype
		message := shared.Filter{TimeStamp: time.Now(), SourceIP: "192.168.102.22"}
		fmt.Print("Event Type:")
		event, _ := reader.ReadString('\n')
		event = strings.TrimSpace(event)

		if event == "exit" {
			fmt.Println("Exiting...")
			break
		}
		message.EventType = event

		//Severity Level
		fmt.Print("Severity Level:")
		severity, _ := reader.ReadString('\n')
		severity = strings.TrimSpace(severity)

		if severity == "exit" {
			fmt.Println("Exiting...")
			break
		}
		severityInt, err := strconv.Atoi(severity)
		if err != nil {
			fmt.Println("Error parsing int from console, pls input severity as integer")
			continue
		}

		message.Severity = severityInt

		//UserId
		fmt.Print("UserId:")
		userID, _ := reader.ReadString('\n')
		userID = strings.TrimSpace(userID)

		if userID == "exit" {
			fmt.Println("Exiting...")
			break
		}
		message.UserID = userID

		var response string
		if err := cli.CallResult(ctx, "Filter", []shared.Filter{message}, &response); err != nil {
			log.Printf("Call failed: %v", err)
			continue
		}

		fmt.Printf("Server response: %s\n", response)
	}

	cli.Close()
}
