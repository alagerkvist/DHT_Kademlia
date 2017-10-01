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
	myRoutingTable *RoutingTable
	fileManager *FileManager
}


func (network *Network) Listen() {

	fmt.Println("Listen " + network.myRoutingTable.me.String())

	ipAndPort := strings.Split(network.myRoutingTable.me.Address, ":")
	port, err := strconv.Atoi(ipAndPort[1])

	p := make([]byte, 4096)

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

		n,remoteaddr,err := ser.ReadFromUDP(p)

		unMarshalMessage := &ProtocolPackage{}
		err = proto.Unmarshal(p[:n], unMarshalMessage)
		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}

		//new contact and add it to bucket
		newContact := &Contact{
			ID: NewKademliaIDFromBytes(unMarshalMessage.FindID),
			Address: *unMarshalMessage.Address,
		}
		newContact.CalcDistance(network.myRoutingTable.me.ID)
		network.myRoutingTable.createTask(addContact, nil, newContact)

		//fmt.Println(unMarshalMessage.GetMessageSent())

		switch unMarshalMessage.GetMessageSent() {
		case ProtocolPackage_PING:
			//fmt.Printf("Ping")
			network.processPing(unMarshalMessage, remoteaddr, ser)
			break;
		case ProtocolPackage_STORE:
			//fmt.Printf("store")
			network.processStoreMessage(unMarshalMessage, remoteaddr, ser)
			break;
		case ProtocolPackage_FINDNODE:
			//fmt.Println("\n\n --- find node --- \n")
			network.processFindConctactMessage(unMarshalMessage, remoteaddr, ser)

			break;
		case ProtocolPackage_FINDVALUE:
			//TODO process FindValue
			//fmt.Printf("find value")
			break;

		}

	}

}


func (network *Network) processPing(protocolPackage *ProtocolPackage, remoteaddr *net.UDPAddr, ser *net.UDPConn){
	fmt.Print("Ping processor:   ")
	fmt.Print(remoteaddr)
	fmt.Print(protocolPackage)

	typeOfMessage := ProtocolPackage_PING
	pongPacket := &ProtocolPackage{
		ClientID: network.myRoutingTable.me.ID.getBytes(),
		Address: &network.myRoutingTable.me.Address,
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

	kclosetContacts := network.myRoutingTable.FindClosestContacts(NewKademliaIDFromBytes(protocolPackage.FindID), bucketSize)

	sendContacts := make([]*ProtocolPackage_ContactInfo, 0)

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
		Address: &network.myRoutingTable.me.Address,
		MessageSent: &typeOfMessage,
		ContactsKNearest: sendContacts,
	}
	marshalledNodesPacket, err := proto.Marshal(responsePkg)
	if err == nil {
		fmt.Println("\n\n Response will be send from: " + network.myRoutingTable.me.String() + " about closest nodes\n")
		_,err := ser.WriteToUDP(marshalledNodesPacket, remoteaddr)
		if err != nil {
			fmt.Printf("Couldn't send response %v", err)
		}

	}else {
		log.Fatal("marshaling find contact response error: ", err)
	}
}


func (network *Network) processStoreMessage(protocolPackage *ProtocolPackage, remoteaddr *net.UDPAddr, ser *net.UDPConn){
	var id *string = protocolPackage.StoredeID
	var base64File *string = protocolPackage.File
	network.fileManager.checkAndStore(*id, *base64File)
}



func (network *Network) Sender (marshaledObject []byte, address string, answerWanted bool) (*ProtocolPackage){

	p :=  make([]byte, 2048)
	conn, err := net.Dial("udp", address)

	if err != nil {
		fmt.Printf("126 Some error %v", err)
		return nil
	}
	fmt.Fprintf(conn, string(marshaledObject))
	if answerWanted {
		n, err := conn.Read(p)

		if err == nil {

			unMarshalledResponse := &ProtocolPackage{}
			err = proto.Unmarshal(p[:n], unMarshalledResponse)

			//new contact and add it to bucket
			fmt.Println(unMarshalledResponse)
			newContact := &Contact{
				ID:      NewKademliaIDFromBytes(unMarshalledResponse.ClientID),
				Address: address,
			}

			newContact.CalcDistance(network.myRoutingTable.me.ID)
			network.myRoutingTable.createTask(addContact, nil, newContact)

			switch unMarshalledResponse.GetMessageSent() {
			case ProtocolPackage_PING:
				//fmt.Printf("Ping")
				break;
			case ProtocolPackage_STORE:
				//fmt.Printf("store")
				break;
			case ProtocolPackage_FINDNODE:
				//fmt.Printf("find node")
				break;
			case ProtocolPackage_FINDVALUE:
				//fmt.Printf("find value")
				break;

			}

			if err != nil {
				log.Fatal("unmarshaling error: ", err)
			}

			conn.Close()
			return unMarshalledResponse
		} else {
			fmt.Printf("175 Some error %v\n", err)
		}
		conn.Close()
	}
	return nil
}




func SendFindDataMessage()  {

}




func (network *Network) SendPingMessage(contact *Contact) {
	network.marshalPing(contact)
	//fmt.Println(result.Address)
}

func (network *Network) marshalPing(contacts *Contact) (*ProtocolPackage) {
	typeOfMessage := ProtocolPackage_PING
	marshalPackage := &ProtocolPackage{
		ClientID: network.myRoutingTable.me.ID.getBytes(),
		Address: proto.String(network.myRoutingTable.me.Address),
		MessageSent: &typeOfMessage,
	}

	data, err := proto.Marshal(marshalPackage)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	return network.Sender(data, contacts.Address, true)
}



func (network *Network) SendFindContactMessage(contact *Contact, findThisID *KademliaID) *ContactCandidates{

	result := network.marshalFindContact(findThisID, contact)
	//fmt.Println(len(result.ContactsKNearest))

	newCandidates := &ContactCandidates{}
	newContacts := make([]Contact, len(result.ContactsKNearest))

	//fmt.Println("\n--- Contacts recus: ")
	for i := 0 ; i < len(result.ContactsKNearest) ; i++{
		//Create the new contact
		newContacts[i] = NewContact(NewKademliaIDFromBytes(result.ContactsKNearest[i].ContactID), *result.ContactsKNearest[i].Address)
		newContacts[i].distance = NewKademliaIDFromBytes(result.ContactsKNearest[i].Distance)
		//fmt.Println(newContacts[i].String())
	}
	//fmt.Println("---\n")

	newCandidates.Append(newContacts)

	return newCandidates
}

func (network *Network) marshalFindContact(findThisID *KademliaID, contact *Contact) (*ProtocolPackage){

	// fmt.Println("Marshal: target " + findThisID.String() + "  ask to " + contact.String() + "\n")

	typeOfMessage := ProtocolPackage_FINDNODE
	fmt.Println(network.myRoutingTable.me)
	marshalPackage := &ProtocolPackage{
		ClientID: network.myRoutingTable.me.ID.getBytes(),
		Address: proto.String(network.myRoutingTable.me.Address),
		MessageSent: &typeOfMessage,
		FindID: findThisID.getBytes(),
	}

	data, err := proto.Marshal(marshalPackage)

	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	return network.Sender(data, contact.Address, true)
}



func (network *Network) SendFindDataMessage(hash string) {
	// TODO
	// Serialize

}

func (network *Network) SendStoreMessage(fileName string, data string, contactsToSend []NodeToCheck){

	for i := 0 ; i < len(contactsToSend) ; i++{
		go network.marshalStore(fileName, data, contactsToSend[i].contact)
	}

}


func (network *Network) marshalStore(fileName string, data string, contact *Contact) (*ProtocolPackage){

	typeOfMessage := ProtocolPackage_STORE

	marshalPackage := &ProtocolPackage{
		ClientID: network.myRoutingTable.me.ID.getBytes(),
		Address: proto.String(network.myRoutingTable.me.Address),
		MessageSent: &typeOfMessage,
		StoredeID: &fileName,
		File: &data,
	}

	marshalData, err := proto.Marshal(marshalPackage)

	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	return network.Sender(marshalData, contact.Address, false)
}





func (network *Network) GetMyRoutingTable() *RoutingTable{
	return network.myRoutingTable
}


func (network *Network) PrintNetwork () {
	fmt.Println(network.myRoutingTable.me.ID)
	//fmt.Println(network.myContact.distance)
	fmt.Println(network.myRoutingTable.me.Address)
	return
}


func (network *Network) TestKademliaPing(contact *Contact) {
	for{
		time.Sleep(2 * time.Second)
		go network.SendPingMessage(contact)
	}
}