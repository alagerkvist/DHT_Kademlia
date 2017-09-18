package kademlia

import (
	"net"
	"fmt"
	//"bufio"
	//"log"
	//"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"log"
	"strings"
	"strconv"
	"time"
)

type Network struct {
	myContact *Contact
	myRoutingTable *RoutingTable
}


func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_,err := conn.WriteToUDP([]byte("From server: Hello I got your mesage "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func (network *Network) Listen() {

	// TODO
	//socket listening different events
	// PingMessage
	// FindContactMessage
	// FindDataMessage
	// StoreMessage
	//fmt.Println("****Listener***")
	//fmt.Println(network.myContact.ID)
	//fmt.Println(network.myContact.Address)
	ipAndPort := strings.Split(network.myContact.Address, ":")
	port, err := strconv.Atoi(ipAndPort[1])

	p := make([]byte, 2048)

	addr := net.UDPAddr{
		Port: port,
		IP: net.ParseIP(ipAndPort[0]),
	}
	//fmt.Println("addr")
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error A %v\n", err)
		return
	}
	for {

		//data,remoteaddr,err := ser.ReadFromUDP(p)
		//fmt.Println("listen")

		n,remoteaddr,err := ser.ReadFromUDP(p)

		fmt.Println("¿¿====")
		unMarshalMessage := &ProtocolPackage{}
		err = proto.Unmarshal(p[:n], unMarshalMessage)
		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}

		fmt.Println(unMarshalMessage.GetMessageSent())
		switch unMarshalMessage.GetMessageSent() {
		case ProtocolPackage_PING:
			fmt.Printf("Ping")
			network.processPing(unMarshalMessage, remoteaddr, ser)
			break;
		case ProtocolPackage_STORE:
			//TODO process Store
			fmt.Printf("store")
			break;
		case ProtocolPackage_FINDNODE:
			fmt.Println("&&&&&&&&&&&&&&&find node")
			network.processFindConctactMessage(unMarshalMessage, remoteaddr, ser)

			break;
		case ProtocolPackage_FINDVALUE:
			//TODO process FindValue
			fmt.Printf("find value")
			break;

		}


		/*
		//_,remoteaddr,err := ser.ReadFromUDP(p)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if err !=  nil {
			fmt.Printf("Some error B  %v", err)
			continue
		}*/
		//go unmarshallData(data)
		//go sendResponse(ser, remoteaddr)
	}

	//unserialize
}

func (network *Network) GetMyContact() *Contact{
	return network.myContact
}

func (network *Network) GetMyRoutingTable() *RoutingTable{
	return network.myRoutingTable
}

func unmarshallData(data int) {
	//newTest := &ProtocolPackage{}
	//err = proto.Unmarshal(data, newTest)
	//if err != nil {
	//	log.Fatal("unmarshaling error: ", err)
	//}
}

func (network *Network) Sender (marshaledObject []byte, address string) (*ProtocolPackage){

	fmt.Println("4444sender")
	fmt.Println(address)
	p :=  make([]byte, 2048)
	conn, err := net.Dial("udp", address)
	// //net.Dial("udp", "127.0.0.1:8080")

	if err != nil {
		fmt.Printf("126 Some error %v", err)
		return nil
	}
	fmt.Fprintf(conn, string(marshaledObject))
	fmt.Println(p)
	n,err := conn.Read(p)
	//fmt.Println(p)
	//a, err := bufio.NewReader(conn).Read(p)
	//fmt.Println(a)
	fmt.Println("SEnder138")
	if err == nil {
		fmt.Println("Response from the server")
		unMarshalledResponse := &ProtocolPackage{}
		err = proto.Unmarshal(p[:n], unMarshalledResponse)

		//new contact and add it to bucket
		/*newContact := &Contact{
			ID: NewKademliaIDFromBytes(unMarshalledResponse.FindID),
			Address: address,
		}*/

		//newContact.CalcDistance(network.myContact.ID)
		//network.myRoutingTable.AddContact(*newContact)


		switch unMarshalledResponse.GetMessageSent() {
		case ProtocolPackage_PING:
			fmt.Printf("Ping")
			break;
		case ProtocolPackage_STORE:
			fmt.Printf("store")
			break;
		case ProtocolPackage_FINDNODE:
			fmt.Printf("find node")
			break;
		case ProtocolPackage_FINDVALUE:
			fmt.Printf("find value")
			break;

		}


		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}

		fmt.Println("Unmarshalled to: %+v", unMarshalledResponse)
		conn.Close()
		return unMarshalledResponse
	} else {
		fmt.Printf("175 Some error %v\n", err)
	}
	conn.Close()


	/*testingRep := make([]*ProtocolPackage_ContactInfo, 0)
	//testingRep := [3]kademlia.ProtocolPackage_ContactInfo{}
	t := ProtocolPackage_ContactInfo{
		ContactID: []byte("Client!!!!!"),
		Address: proto.String("localhost"),
		Distance: []byte("FAR!"),
	}
*/
	//testingRep = append(testingRep, &t)
	return nil
}

func processReceivedMessage () {

}

func (network *Network) processPing(protocolPackage *ProtocolPackage, remoteaddr *net.UDPAddr, ser *net.UDPConn){
	fmt.Print("Ping processor:   ")
	fmt.Print(remoteaddr)
	fmt.Print(protocolPackage)

	typeOfMessage := ProtocolPackage_PING
	pongPacket := &ProtocolPackage{
		ClientID: network.myContact.ID.getBytes(),
		Address: &network.myContact.Address,
		MessageSent: &typeOfMessage,
	}
	marshalledpongPacket, err := proto.Marshal(pongPacket)
	if err == nil {
		_,err := ser.WriteToUDP(marshalledpongPacket, remoteaddr)
		if err != nil {
			fmt.Printf("Couldn’t send response %v", err)}
	}else {
		log.Fatal("marshaling pong error: ", err)
	}

}

func (network *Network) processFindConctactMessage(protocolPackage *ProtocolPackage, remoteaddr *net.UDPAddr, ser *net.UDPConn)  {
	fmt.Println("processFindConctactMessage procesor")
	fmt.Println(protocolPackage)
	kclosetContacts := network.myRoutingTable.FindClosestContacts(NewKademliaIDFromBytes(protocolPackage.FindID), bucketSize)

	sendContacts := make([]*ProtocolPackage_ContactInfo, 0)
	//testingRep := [3]kademlia.ProtocolPackage_ContactInfo{}

	for i :=0; i < len(kclosetContacts); i++{
		contact := ProtocolPackage_ContactInfo{
			ContactID: kclosetContacts[i].ID.getBytes(),
			Address: &kclosetContacts[i].Address,
			Distance: kclosetContacts[i].distance.getBytes(),
		}
		sendContacts = append(sendContacts, &contact)
	}
	typeOfMessage := ProtocolPackage_FINDNODE
	responsePkg := &ProtocolPackage{
		Address: &network.myContact.Address,
		MessageSent: &typeOfMessage,
		ContactsKNearest: sendContacts,
	}
	marshalledNodesPacket, err := proto.Marshal(responsePkg)
	if err == nil {
		fmt.Println("¿¿¿¿¿¿¿¿¿¿¿¿¿¿¿¿¿")
		fmt.Println(marshalledNodesPacket)


		_,err := ser.WriteToUDP(marshalledNodesPacket, remoteaddr)
		if err != nil {
			fmt.Printf("Couldn't send response %v", err)
		}

	}else {
		fmt.Println("??????????????")

		log.Fatal("marshaling find contact rsponse error: ", err)
	}

}

func SendFindDataMessage()  {

}

func SendStoreMessage()  {

}

func (network *Network) SendPingMessage(contact *Contact) {
	network.marshalPing(contact)
	//fmt.Println(result.Address)
}


func (network *Network) marshalPing(contacts *Contact) (*ProtocolPackage) {
	typeOfMessage := ProtocolPackage_PING
	marshalPackage := &ProtocolPackage{
		ClientID: network.myContact.ID.getBytes(),
		Address: proto.String(network.myContact.Address),
		MessageSent: &typeOfMessage,
	}

	data, err := proto.Marshal(marshalPackage)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	return network.Sender(data, contacts.Address)
}

func (network *Network) marshalFindContact(findThisID *KademliaID, contacts *Contact) (*ProtocolPackage){
	typeOfMessage := ProtocolPackage_FINDNODE

	fmt.Println(findThisID)
	fmt.Println(contacts)
	fmt.Println(findThisID.getBytes())
	marshalPackage := &ProtocolPackage{
		ClientID: network.myContact.ID.getBytes(),
		Address: proto.String(network.myContact.Address),
		MessageSent: &typeOfMessage,
		FindID: findThisID.getBytes(),
	}

	fmt.Println(marshalPackage)
	data, err := proto.Marshal(marshalPackage)
	fmt.Println(data)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	return network.Sender(data, contacts.Address)
}


func (network *Network) SendFindContactMessage(contact *Contact, findThisID *KademliaID) (*ContactCandidates){
	// TODO
	// Serialize
	//
	result := network.marshalFindContact(findThisID, contact)
	fmt.Println("****")
	fmt.Println(result)
	return &ContactCandidates{}
	/*
	var contactsReceived = []ContactCandidates
	if len != nil {
		contactsReceived = ContactCandidates{make([]Contact, 0, len(result.ContactsKNearest))}
	}

	for i:=0 ; i < len(result.ContactsKNearest) ; i++{
		///var kademliaId *KademliaID = result.ContactsKNearest[i].ContactID
		newContact := Contact{ID: NewKademliaIDFromBytes(result.ContactsKNearest[i].ContactID), Address: *result.ContactsKNearest[i].Address}
		contactsReceived.contacts = append(contactsReceived.contacts, newContact)
	}
	return &contactsReceived*/
}


func (network *Network) SendFindDataMessage(hash string) {
	// TODO
	// Serialize

}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
	// Serialize
}


func (network *Network) PrintNetwork () {
	//fmt.Println(network.myContact.ID)
	//fmt.Println(network.myContact.distance)
	fmt.Println(network.myContact.Address)
	//fmt.Println(network.myRoutingTable)
	return
}


func (network *Network) TestKademliaPing(contact *Contact) {
	for{
		time.Sleep(2 * time.Second)
		go network.SendPingMessage(contact)
	}
}