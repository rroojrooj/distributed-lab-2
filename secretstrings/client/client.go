package main

import (
	"flag"
	"fmt"
	"net/rpc"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

func main() {
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	fmt.Println("Server: ", *server)

	client, err := rpc.Dial("tcp", *server) // dial the server address that the client passed
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer client.Close()

	// Ask the user for input
	fmt.Print("Enter a string to reverse: ")
	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	request := stubs.Request{Message: input} // Define a request of type 'stubs Request Message' with the value from "input"
	response := new(stubs.Response)
	err = client.Call(stubs.PremiumReverseHandler, request, response) // call the function in server.go, in this case is the PremiumReversHandlerer
	if err != nil {
		fmt.Println("Error calling RPC:", err)
		return
	}

	fmt.Println("Reversed string:", response.Message) // print the message that get out in our response
}
