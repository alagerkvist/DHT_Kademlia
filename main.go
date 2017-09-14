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
	typeOfMessage := kademlia.ProtocolPackage_PING
	var testingRep = make([][]byte, 3)
	//testingRep[0] = make([]byte, 3)
	testingRep[0] = []byte("Client!!!!!")
	testingRep[1] = []byte("IPPP!!!!!")
	testingRep[2] = []byte("PORT!!!!!")

	//testingRep[1] = make([]byte, 3)
	//testingRep[3] = []byte("Client1111")
	/*
	testingRep[1][1] = byte("IPPP1111")
	testingRep[1][2] = byte("PORT1111")
	testingRep[2] = make(chan []byte, 3)
	testingRep[2][0] = byte("Client222")
	testingRep[2][1] = byte("IPPP222")
	testingRep[2][2] = byte("PORT222")
	*/
	//p := &testWhat
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