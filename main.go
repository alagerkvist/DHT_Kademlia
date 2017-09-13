package main

import (
	"./kademlia"
	//"time"
	"net"
	"fmt"
	"bufio"
	"github.com/golang/protobuf/proto"
	"log"
)

func main() {
	//kademlia.Listen("127.0.0.1", 8080);
	testWhat := kademlia.ProtocolPackage_PING
	//p := &testWhat
	testProto := &kademlia.ProtocolPackage{
		ClientID: []byte("String"),
		Ip: proto.String("localhost"),
		ListenPort: proto.Int32(1234),
		MessageSent: &testWhat,
	}

	data, err := proto.Marshal(testProto)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	newTest := &kademlia.ProtocolPackage{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	// Now test and newTest contain the same data.
	if testProto.GetIp() != newTest.GetIp() {
		log.Fatalf("data mismatch %q != %q", testProto.GetIp(), newTest.GetIp())
	}

	log.Printf("Unmarshalled to: %+v", newTest)
/*
	for {
		time.Sleep(2 * time.Second)
		//go SendPingMessageFake()
	}
*/
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