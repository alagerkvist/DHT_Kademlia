package kademlia

import (
	"net"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"strings"
	"strconv"
	"encoding/base64"
	"math"
	"time"
)

type Network struct {
	myRoutingTable *RoutingTable
	FileManager *FileManager
}

const packetSize = 8096

const PRINT_PONG = true

/** Listen
*	PARAM: network
* Listen for incoming requests
 */
func (network *Network) Listen() {

	fmt.Println("Listen " + network.myRoutingTable.me.String())

	ipAndPort := strings.Split(network.myRoutingTable.me.Address, ":")
	port, err := strconv.Atoi(ipAndPort[1])

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
		p := make([]byte, packetSize)

		n,remoteaddr,err := ser.ReadFromUDP(p)
		//fmt.Println(p)
		unMarshalMessage := &ProtocolPackage{}
		err = proto.Unmarshal(p[:n], unMarshalMessage)
		if err != nil {
			log.Fatal("54: unmarshaling error: ", err)
		}
		//new contact and add it to bucket

		newContact := &Contact{
			ID: NewKademliaIDFromBytes(unMarshalMessage.ClientID),
			Address: *unMarshalMessage.Address,
		}
		newContact.CalcDistance(network.myRoutingTable.me.ID)
		network.myRoutingTable.createTask(addContact, nil, newContact)

/*		fmt.Println("*******************************************************************")
		fmt.Println("*******************************************************************")
		fmt.Println(unMarshalMessage.GetMessageSent())
		fmt.Println("*******************************************************************")
		fmt.Println("*******************************************************************")*/

		switch unMarshalMessage.GetMessageSent() {
		case ProtocolPackage_PING:
			//fmt.Printf("Ping")
			//network.processPing(unMarshalMessage, remoteaddr, ser)
			//fmt.Printf("Ping")
			go network.processPing(unMarshalMessage, remoteaddr, ser)
			break;
		case ProtocolPackage_STORE:
			//fmt.Println("*******************            store       Listerner     *******************")
			go network.processStoreMessage(unMarshalMessage, remoteaddr, ser)
			break;
		case ProtocolPackage_FINDNODE:
			//fmt.Println("\n\n --- find node --- \n")
			//network.processFindConctactMessage(unMarshalMessage, remoteaddr, ser)
			//fmt.Println("\n\n --- find node --- \n")
			go network.processFindConctactMessage(unMarshalMessage, remoteaddr, ser)

			break;
		case ProtocolPackage_FINDVALUE:
			go network.processFindValue(unMarshalMessage, remoteaddr, ser)

			break;

		}

	}

}

/** processPing
*	PARAM: network
			protocolPackage: the message receive
			remoteaddr: udp addres
			ser: udp connection
* Process the ping request
 */
func (network *Network) processPing(protocolPackage *ProtocolPackage, remoteaddr *net.UDPAddr, ser *net.UDPConn){


	typeOfMessage := ProtocolPackage_PING
	pongPacket := &ProtocolPackage{
		ClientID: network.myRoutingTable.me.ID.getBytes(),
		Address: &network.myRoutingTable.me.Address,
		MessageSent: &typeOfMessage,
	}
	marshalledpongPacket, err := proto.Marshal(pongPacket)
	if err == nil {
		if(PRINT_PONG){
			fmt.Print("Ping processor:   ")
			fmt.Print(remoteaddr)
			fmt.Print(protocolPackage)
			fmt.Println("pong")
		}

		_,err := ser.WriteToUDP(marshalledpongPacket, remoteaddr)
		if err != nil {
			//fmt.Printf("Couldn’t send response %v", err)
			network.myRoutingTable.RemoveContact(NewContact(NewKademliaIDFromBytes(protocolPackage.ClientID), *protocolPackage.Address))
		}

	}else {
		log.Fatal("marshaling pong error: ", err)
	}

}

/** processFindContactMessage
*	PARAM: network
			protocolPackage: the message receive
			remoteaddr: udp addres
			ser: udp connection
* Process the find contact request
 */
func (network *Network) processFindConctactMessage(protocolPackage *ProtocolPackage, remoteaddr *net.UDPAddr, ser *net.UDPConn)  {
	respChan := make(chan []Contact)
	network.myRoutingTable.createTask(getClosest, respChan, &Contact{NewKademliaIDFromBytes(protocolPackage.FindID), "", nil})
	kclosetContacts := <- respChan

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
		ClientID: network.myRoutingTable.me.ID.getBytes(),
		Address: &network.myRoutingTable.me.Address,
		MessageSent: &typeOfMessage,
		ContactsKNearest: sendContacts,
	}
	marshalledNodesPacket, err := proto.Marshal(responsePkg)
	if err == nil {
		//fmt.Println("Response will be send from: " + network.myRoutingTable.me.String() + " about closest nodes")
		_,err := ser.WriteToUDP(marshalledNodesPacket, remoteaddr)
		if err != nil {
			//fmt.Printf("Couldn't send response %v", err)
			network.myRoutingTable.RemoveContact(NewContact(NewKademliaIDFromBytes(protocolPackage.ClientID), *protocolPackage.Address))
		}

	}else {
		log.Fatal("marshaling find contact response error: ", err)
	}
}

/** processStoreMessage
*	PARAM: network
			protocolPackage: the message receive
			remoteaddr: udp addres
			ser: udp connection
* Process the store message
*/
func (network *Network) processStoreMessage(protocolPackage *ProtocolPackage, remoteaddr *net.UDPAddr, ser *net.UDPConn){
	var id *string = protocolPackage.StoredeID
	var base64File *string = protocolPackage.File
	if network.FileManager.checkIfFileExist(filesDirectory + *id){
		network.FileManager.updateTime(*id)
	} else{
		network.FileManager.CheckAndStore(*id, *base64File)
		network.SetExpirationTime(*id)
	}
}


/** processFindValue
*	PARAM: network
			protocolPackage: the message receive
			remoteaddr: udp addres
			ser: udp connection
* Process the find value request
 */
func (network *Network) processFindValue(protocolPackage *ProtocolPackage, remoteaddr *net.UDPAddr, ser *net.UDPConn){
	var id string = NewKademliaIDFromBytes(protocolPackage.FindValue).String()

	if !network.FileManager.checkIfFileExist(filesDirectory + id){
		network.processFindConctactMessage(protocolPackage, remoteaddr, ser)
	} else {
		typeOfMessage := ProtocolPackage_FINDVALUE

		file := base64.StdEncoding.EncodeToString(network.FileManager.readData(filesDirectory + id))

		responsePkg := &ProtocolPackage{
			ClientID: network.myRoutingTable.me.ID.getBytes(),
			Address: &network.myRoutingTable.me.Address,
			MessageSent: &typeOfMessage,
			ContactsKNearest: nil,
			File: &file,
		}

		marshalledNodesPacket, err := proto.Marshal(responsePkg)
		if err == nil {
			fmt.Println("Response will be send from: " + network.myRoutingTable.me.String() + " about a file")
			_,err := ser.WriteToUDP(marshalledNodesPacket, remoteaddr)
			if err != nil {
				fmt.Printf("Couldn't send response %v", err)
				network.myRoutingTable.RemoveContact(NewContact(NewKademliaIDFromBytes(protocolPackage.ClientID), *protocolPackage.Address))
			}

		}else {
			log.Fatal("marshaling find contact response error: ", err)
		}
	}
}


/** Sender
*	PARAM: network
			marshaledObject: the message to send
			address: udp addres to send
			answerWanted: if answer is expected
* Process the ping request
 */
func (network *Network) Sender(marshaledObject []byte, address string, answerWanted bool) (*ProtocolPackage){

	fmt.Println("sender", address)
	p :=  make([]byte, packetSize)
	//fmt.Println(marshaledObject)
	//fmt.Println(string(marshaledObject))
	conn, err := net.DialTimeout("udp", address, time.Second * 2)

	if err != nil {
		fmt.Printf("126 Some error %v", err)
		return nil
	}
	//fmt.Fprintf(conn, string(marshaledObject))
	conn.SetWriteDeadline(time.Now().Add(time.Second * 2))
	ansWrite , errorWrite := conn.Write(marshaledObject)

	if errorWrite != nil {
		return nil
	}

	fmt.Println(ansWrite)

	if answerWanted {
		conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		n, errRead := conn.Read(p)
		if errRead != nil {
			return nil
		}

		if err == nil && errRead == nil {

			unMarshalledResponse := &ProtocolPackage{}
			err = proto.Unmarshal(p[:n], unMarshalledResponse)


			//new contact and add it to bucket
			//fmt.Println(unMarshalledResponse)
			newContact := &Contact{
				ID:      NewKademliaIDFromBytes(unMarshalledResponse.ClientID),
				Address: address,
			}

			newContact.CalcDistance(network.myRoutingTable.me.ID)
			network.myRoutingTable.createTask(addContact, nil, newContact)

			if err != nil {
				log.Fatal("298: unmarshaling error: ", err)
			}

			conn.Close()
			//fmt.Println("unMarshalledResponse")
			//fmt.Println(unMarshalledResponse)
			return unMarshalledResponse
		} else {
			fmt.Printf("175 Some error %v\n", err)
		}
	}
	conn.Close()
	return nil
}


/** sendFindDataValue
*	PARAM: network
			id: id of the file to find
			contact: the contact to send the message
 */
func (network *Network) SendFindDataValue(id KademliaID, contact *Contact) Response{

	result := network.marshalFindValue(id, contact)

	if result == nil{
		return Response{nil, nil, contact, true}
	}


	if result.ContactsKNearest != nil{
		newCandidates := &ContactCandidates{}
		newContacts := make([]Contact, len(result.ContactsKNearest))

		for i := 0 ; i < len(result.ContactsKNearest) ; i++{
			//Create the new contact
			newContacts[i] = NewContact(NewKademliaIDFromBytes(result.ContactsKNearest[i].ContactID), *result.ContactsKNearest[i].Address)
			newContacts[i].CalcDistance(&id)
			//fmt.Println(newContacts[i].String())
		}

		newCandidates.Append(newContacts)
		return Response{newCandidates, nil, nil, false}
	}

	return Response{nil, result.File, contact, false}
}


/**marshalFindValue
*	PARAM: network
			id: id of the file to find
			contact: the contact to send the message
* Prepare the request to send
 */
func (network *Network) marshalFindValue(id KademliaID, contact *Contact) (*ProtocolPackage) {
	typeOfMessage := ProtocolPackage_FINDVALUE

	marshalPackage := &ProtocolPackage{
		ClientID: network.myRoutingTable.me.ID.getBytes(),
		Address: proto.String(network.myRoutingTable.me.Address),
		MessageSent: &typeOfMessage,
		FindValue: id.getBytes(),
	}

	data, err := proto.Marshal(marshalPackage)

	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	return network.Sender(data, contact.Address, true)
}


/** SendPingMessage
*	PARAM: network
			contact: the contact to send the message
*  Prepare the request to send
 */
func (network *Network) SendPingMessage(contact *Contact) {
	network.marshalPing(contact)
}


/** marshalPing
*	PARAM: network
			contact: the contact to send the message
* Prepare packet for ping
 */
func (network *Network) marshalPing(contacts *Contact) (*ProtocolPackage) {
	fmt.Println("marshallping")
	typeOfMessage := ProtocolPackage_PING

	marshalPackage := &ProtocolPackage{
		ClientID: network.myRoutingTable.me.ID.getBytes(),
		Address: proto.String(network.myRoutingTable.me.Address),
		MessageSent: &typeOfMessage,
	}
	//fmt.Println("ok1")
	data, err := proto.Marshal(marshalPackage)
	//fmt.Println(marshalPackage)
	//fmt.Println(data)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	//fmt.Println("ok2")
	return network.Sender(data, contacts.Address, true)
}


/** SendFindContactMessage
*	PARAM: network
			findThisID: id of the contact to find
			contact: the contact to send the message
 */
func (network *Network) SendFindContactMessage(findThisID *KademliaID, contact *Contact) Response{
	result := network.marshalFindContact(findThisID, contact)
	fmt.Println("findContact")
	if result == nil{
		return Response{nil, nil, contact, true}
	}

	newCandidates := &ContactCandidates{}
	newContacts := make([]Contact, len(result.ContactsKNearest))

	for i := 0 ; i < len(result.ContactsKNearest) ; i++{
		//Create the new contact
		newContacts[i] = NewContact(NewKademliaIDFromBytes(result.ContactsKNearest[i].ContactID), *result.ContactsKNearest[i].Address)
		newContacts[i].CalcDistance(findThisID)
	}
	//fmt.Println("---\n")

	newCandidates.Append(newContacts)

	return Response{newCandidates, nil, nil, false}
}


/** marshalFindContact
*	PARAM: network
			findThisID: id of the conttact to find
			contact: the contact to send the message
 */
func (network *Network) marshalFindContact(findThisID *KademliaID, contact *Contact) (*ProtocolPackage){

	// fmt.Println("Marshal: target " + findThisID.String() + "  ask to " + contact.String() + "\n")

	typeOfMessage := ProtocolPackage_FINDNODE
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


/** SendStoreMessage
*	PARAM: network
			fileName: file to store
			id: id of the file to find
			contact: the contact to send the message
 */
func (network *Network) SendStoreMessage(fileName string, data string, contactsToSend []NodeToCheck){
	fmt.Println("==================== SEND STORE   ==================")
	fmt.Println(len(contactsToSend))
	for i := 0 ; i < len(contactsToSend) ; i++{
		network.marshalStore(fileName, data, contactsToSend[i].contact)
	}

}


/** marshalStore
*	PARAM: network
			fileName: id of the file to send
			data: the data to send
			contact: the contact to send
 */
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
	fmt.Println("MARSHALL STORE")
	return network.Sender(marshalData, contact.Address, false)
}


/** GetMyRoutingTable
*	return the routing table of network
 */
func (network *Network) GetMyRoutingTable() *RoutingTable{
	return network.myRoutingTable
}


/** SetExpirationTime
*	PARAM: network
			fileName: filename to put the expiration time
 */
func (network *Network) SetExpirationTime(fileName string){
	fileInfo := network.FileManager.filesStored[fileName]
	responseChannel := make(chan []Contact)
	fileID := NewKademliaID(fileName)
	network.myRoutingTable.createTask(getClosest, responseChannel, &Contact{fileID, "", nil})
	closestContacts := <- responseChannel
	me := network.myRoutingTable.me
	me.CalcDistance(fileID)

	for i:=0 ; i < len(closestContacts) ; i++{
		if i > 13 || me.Less(&closestContacts[i]) || me.ID.Equals(closestContacts[i].ID){
				fileInfo.expirationTime = 24 * math.Exp(-(float64(i)))
				break
		}
	}

}


/** PrintNetwork
* Print the network parameter
 */
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