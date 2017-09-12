package main

import (
	"./kademlia"
	"time"
	"net"
	"fmt"
	"bufio"
)

func main() {
	kademlia.Listen("127.0.0.1", 8080);

	for {
		time.Sleep(2 * time.Second)
		//go SendPingMessageFake()
	}

}

func SendPingMessageFake () {
	// TODO
	p :=  make([]byte, 2048)
	conn, err := net.Dial("udp", "127.0.0.1:8080")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(conn, "Hi UDP Server, How are you doing?")
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()

}