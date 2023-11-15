package main

import (
	"flag"
	"math/rand"
	"net"
	"net/rpc"
	"time"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

/** Super-Secret `reversing a string' method we can't allow clients to see. **/
func ReverseString(s string, i int) string {
	time.Sleep(time.Duration(rand.Intn(i)) * time.Second)
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

type SecretStringOperation struct{}

func (s *SecretStringOperation) Reverse(req stubs.Request, res *stubs.Response) (err error) {
	res.Message = ReverseString(req.Message, 10)
	return
}

func (s *SecretStringOperation) FastReverse(req stubs.Request, res *stubs.Response) (err error) {
	res.Message = ReverseString(req.Message, 2) //decrease the delay that get passed to ReverseString()
	// res.Message is the RESPONSE para stubs. req.message is the REQUEST para stubs
	return
}

func main() { //Listen for indications from the client on calling the func (s *SecretStringOperation), and will handle the communications by the listener.
	pAddr := flag.String("port", "8030", "Port to listen on") // Server accepts some configs
	flag.Parse()
	rand.Seed(time.Now().UnixNano()) // for reverse string implement
	rpc.Register(&SecretStringOperation{})
	listener, _ := net.Listen("tcp", ":"+*pAddr) // create a listener, listens fos commu in TCP port
	defer listener.Close()
	rpc.Accept(listener)
}
