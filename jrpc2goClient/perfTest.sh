#!/bin/bash

# Path to the compiled Go client application
CLIENT_PATH="./main"  # Update this with the correct path if needed

# Function to send messages from a client
send_messages() {
    (
        for i in {1..100}; do
            echo "Hello $i from client $1"
        done
        echo "exit"  # Send exit command to close the client
    ) | $CLIENT_PATH
    
    if [ $? -ne 0 ]; then
        echo "Client $1 encountered an error."
        return 1
    fi
}

# Start 10 clients
for client_id in {1..10}; do
    send_messages "$client_id" &
    if [ $? -ne 0 ]; then
        echo "Failed to start client $client_id."
        exit 1
    fi
done

# Wait for all background jobs to finish
wait

echo "All clients have finished sending messages."