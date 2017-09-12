package kademlia

import (
	"net"
	"fmt"
)

type Network struct {
}

func Listen(ip string, port int) {
	// TODO
	//socket listening different events
	// PingMessage
	// FindContactMessage
	// FindDataMessage
	// StoreMessage

	//unserialize

	fmt.Println("inside")

	var portStr string = string(port)
	protocol := "udp"

	//Build the address
	udpAddr, err := net.ResolveUDPAddr(protocol, portStr)
	if err != nil {
		fmt.Println("Wrong Address")
		return
	}

	udpConn, err := net.ListenUDP(protocol, udpAddr)
	udpConn.Close();
	if err != nil {
		fmt.Println(err)
	}

}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
	// Serialize
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
	// Serialize
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
	// Serialize
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
	// Serialize
}
