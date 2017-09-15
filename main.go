package main

import (
	"./kademlia"
	//"time"
	"net"
	"fmt"
	"bufio"
	"github.com/golang/protobuf/proto"
	"log"
	"time"
)

func main() {
	//kademlia.Listen("127.0.0.1", 8080);
	typeOfMessage := kademlia.ProtocolPackage_PING

	testingRep := make([]*kademlia.ProtocolPackage_ContactInfo, 0)
	//testingRep := [3]kademlia.ProtocolPackage_ContactInfo{}
	t := kademlia.ProtocolPackage_ContactInfo{
		ContactID: []byte("Client!!!!!"),
		Address: proto.String("localhost"),
		Distance: []byte("FAR!"),
	}

	testingRep = append(testingRep, &t)

	tx := kademlia.ProtocolPackage_ContactInfo{
		ContactID: []byte("Client11"),
		Address: proto.String("localhost1"),
		Distance: []byte("FAR1"),
	}
	testingRep = append(testingRep, &tx)
	testProto := &kademlia.ProtocolPackage{
		ClientID: []byte("String"),

		Ip: proto.String("localhost"),
		ListenPort: proto.Int32(1234),
		MessageSent: &typeOfMessage,
		ContactsKNearest: testingRep,

		//ContactsKNearest: testingRep,
		//ContactsKNearest: []byte("Client:123 ip:111"),
		//ContactsKNearest: []byte("Client:1234 ip:222"),
		//Then we can add both findID and findValue, they are optional
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

	for i := 0; i < len(newTest.ContactsKNearest); i++ {
		log.Print("ClientID: ", string(newTest.ContactsKNearest[i].ContactID[:]))
		log.Print("Address: ", *newTest.ContactsKNearest[i].Address)
		log.Print("Distance: ", string(newTest.ContactsKNearest[i].Distance[:]))

	}

	for {
		time.Sleep(2 * time.Second)
		//go SendPingMessageFake()
	}

}

func SendPingMessageFake () {
	// TODO
	p :=  make([]byte, 2048)
	conn, err := net.DialTimeout("udp", "127.0.0.1:1234", 100)
	// //net.Dial("udp", "127.0.0.1:8080")

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