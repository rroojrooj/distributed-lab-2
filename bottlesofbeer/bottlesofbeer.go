package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"time"
)

var nextAddr string
var thisPort string
var bottles int

type Buddy struct{}

type Args struct {
	Bottles int
}

func (b *Buddy) Sing(args *Args, reply *Args) error {
	fmt.Println("Received token with bottles:", args.Bottles)
	fmt.Printf("%d bottles of beer on the wall, %d bottles of beer. Take one down, pass it around...\n", args.Bottles, args.Bottles)
	args.Bottles--
	*reply = *args
	return nil
}

func startServer() {
	fmt.Println("Starting server on port:", thisPort)
	buddy := new(Buddy)
	rpc.Register(buddy)
	l, err := net.Listen("tcp", ":"+thisPort)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func startSong(bottles int) {
	fmt.Println("Starting song with bottles:", bottles)
	time.Sleep(2 * time.Second)
	args := &Args{Bottles: bottles}
	var reply Args
	for {
		client, err := rpc.Dial("tcp", nextAddr)
		if err != nil {
			fmt.Println("Connection error: ", err)
			time.Sleep(time.Second)
			continue
		}
		err = client.Call("Buddy.Sing", args, &reply)
		if err != nil {
			fmt.Println("RPC error: ", err)
			client.Close()
			time.Sleep(time.Second)
			continue
		}
		client.Close()
		if reply.Bottles == 0 {
			fmt.Println("No more bottles of beer on the wall.")
			return
		}
		args.Bottles = reply.Bottles
		time.Sleep(time.Second)
	}
}

func main() {
	flag.StringVar(&thisPort, "this", "8030", "Port for this process to listen on")
	flag.StringVar(&nextAddr, "next", "localhost:8040", "IP:Port string for next member of the round.")
	flag.IntVar(&bottles, "n", 0, "Bottles of Beer (launches song if not 0)")
	flag.Parse()

	go startServer()

	if bottles > 0 {
		startSong(bottles)
	} else {
		select {}
	}
}
