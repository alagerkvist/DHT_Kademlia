package kademlia

import (
	"fmt"
	"strconv"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
)

const alpha = 3


type Kademlia struct {
	network *Network
}

type NodeToCheck struct{
	contact	*Contact
	alreadyChecked	bool 	//usefull if network < k
}


type Request struct{
	contact *Contact
	endWork bool
}

type Response struct{
	newContacts *ContactCandidates
	data *string
	contactedContact *Contact
	error bool
}

//LookupContact is a method of Kademlia to locate some Node
//Params target: it is the finded contact
func (kademlia *Kademlia) LookupContact(targetId *KademliaID) []NodeToCheck{
	return kademlia.Lookup(targetId, true)
}



func (kademlia *Kademlia) LookupData(fileName string) {
	targetID := NewKademliaID(fileName)
	kademlia.Lookup(targetID, false)

}
//LookupContact is a method of KAdemlia to locate some Data
//PArams hash: it is the finded data with the 160 bits hash
func (kademlia *Kademlia) Lookup(targetID *KademliaID, isForNode bool) []NodeToCheck{
	channelToSendRequest := make(chan Request, alpha)
	channelToReceive := make(chan Response)
	nodesToCheck := make([]NodeToCheck, 0, bucketSize)
	countEndThread := 0
	contactWithFiles := []Contact{}

	//Fulfill this array with at most the k nodes from buckets
	responseChannel := make(chan []Contact)
	kademlia.network.myRoutingTable.createTask(lookUpContact, responseChannel, &Contact{targetID, "", nil})
	firstContacts := <- responseChannel

	for i:=0 ; i<len(firstContacts) ; i++ {
		nodesToCheck = append(nodesToCheck, NodeToCheck{&firstContacts[i], false})
	}

	for i := 0 ; i < alpha ; i++ {
		go kademlia.network.workerFindData(channelToSendRequest, *targetID, channelToReceive, isForNode)

		fmt.Println("sender")
		Print(nodesToCheck)
		channelToSendRequest <- Request{nodesToCheck[i].contact, false}
		nodesToCheck[i].alreadyChecked = true
	}

	for {
		newResponse := <- channelToReceive
		if newResponse.error{
			fmt.Println("Unreachable node")
			kademlia.network.myRoutingTable.createTask(removeContact, nil, newResponse.contactedContact)

		} else if !isForNode && newResponse.newContacts == nil{
			fmt.Println("response: " + *newResponse.data)
			contactWithFiles = append(contactWithFiles, *newResponse.contactedContact)
			kademlia.network.FileManager.CheckAndStore(targetID.String(), *newResponse.data)

			for ; countEndThread + 1 < alpha; countEndThread++{
				//Get the last responses and check if they have the file
				newResponse := <- channelToReceive
				fmt.Println(newResponse)
				if newResponse.newContacts == nil{
					contactWithFiles = append(contactWithFiles, *newResponse.contactedContact)
				}
			}
			sendEndWork(channelToSendRequest, alpha)
			//Find the node to send store message
			for i:=0 ; i < len(nodesToCheck) ; i++{
				canSend := false
				for j:= 0 ; j < len(contactWithFiles) ; j++ {
					if nodesToCheck[i].contact.ID.Equals(contactWithFiles[j].ID) {
						canSend = true
					}
				}
				if(canSend){
					kademlia.network.marshalStore(targetID.String(), *newResponse.data, nodesToCheck[i].contact)
					break
				}
			}

			break
		}

		newContacts := newResponse.newContacts

		for i:=0 ; newContacts != nil && i<newContacts.Len() ; i++{
			//fmt.Println("Contact to add: " + newContacts.contacts[i].String())

			for j:=0 ; j<bucketSize ; j++{
				if !newContacts.contacts[i].ID.Equals(kademlia.network.myRoutingTable.me.ID){
					//case of less than k values
					if j >= len(nodesToCheck) {
						nodesToCheck = append(nodesToCheck, NodeToCheck{&newContacts.contacts[i], false})
						break;
					} else if newContacts.contacts[i].ID.Equals(nodesToCheck[j].contact.ID){
						break;
					} else if  newContacts.contacts[i].Less(nodesToCheck[j].contact){
						//Shift value of the array and insert the new one
						copy(nodesToCheck[j+1:], nodesToCheck[j: len(nodesToCheck) - 1])
						nodesToCheck[j].contact = &(newContacts.contacts[i])
						nodesToCheck[j].alreadyChecked = false
						break;
					}
				}
			}
		}

		nextContactToCheck :=  kademlia.network.getNextContactToAsk(nodesToCheck)
		//Check if still some work to do, if it is not wait until all the workers finish
		if nextContactToCheck == nil{
			countEndThread++
			if countEndThread == alpha {
				if(!isForNode) {
					fmt.Println("Impossible to find the file in the network")
				}
				Print(nodesToCheck)

				//Print(nodesToCheck)
				sendEndWork(channelToSendRequest, alpha)
				break
			}
		} else {
			if countEndThread > 0{
				for ; countEndThread > 0 && nextContactToCheck != nil ; countEndThread--{
					nextContactToCheck = kademlia.network.getNextContactToAsk(nodesToCheck)
				}
			} else{
				channelToSendRequest <- Request{nextContactToCheck, false}
			}
		}
	}

	return nodesToCheck
}


func sendEndWork(channelToSendRequest chan Request, nb int){
	for i := 0 ; i < alpha ; i++{
		channelToSendRequest <- Request{nil, true}
	}
}


func(network *Network) workerFindData(requestsChannel chan Request, targetId KademliaID, responseChannel chan Response, isForNode bool) {

	for {
		request := <- requestsChannel
		if(request.endWork){
			break
		}
		fmt.Print("request: ")
		fmt.Println(request.contact)
		if !isForNode {
			responseChannel <- network.SendFindDataValue(targetId, request.contact)
		} else {
			responseChannel <- network.SendFindContactMessage(&targetId, request.contact)
		}
	}
}


func (network *Network) getNextContactToAsk(nodesToCheck []NodeToCheck) *Contact{
	for i:=0 ; i < len(nodesToCheck) ; i++ {
		if !nodesToCheck[i].alreadyChecked && !nodesToCheck[i].contact.ID.Equals(network.myRoutingTable.me.ID) {
			nodesToCheck[i].alreadyChecked = true
			return nodesToCheck[i].contact
		}
	}
	return nil
}


//Store is the method of KAdemlia to Store data
// Params: data array of Bytes.
func (kademlia *Kademlia) Store(fileName string) {
	fileManager := kademlia.network.FileManager

	//fmt.Println("$$$$$$$$$       saving in other node       $$$$$$$$$")

	if !fileManager.checkIfFileExist(fileName){
		fmt.Println("File not found")
	}else {
		data := fileManager.readData(fileName)
		base64Data := base64.StdEncoding.EncodeToString(data[:])

		//Generate a hash for the name of the file
		hash := sha256.Sum256(data)
		idFile := NewKademliaIDFromBytes(hash[:IDLength])

		fileManager.CheckAndStore(idFile.String(), base64Data)

		contactToSend := kademlia.LookupContact(idFile)
		//fmt.Println("File will be send to these contacts:")
		Print(contactToSend)
		kademlia.network.SendStoreMessage(idFile.String(), base64Data, contactToSend)
	}
}

func (kademlia *Kademlia) PrintFile(fileName string) {
	fileManager := kademlia.network.FileManager
	completeFileName := filesDirectory + fileName
	if !fileManager.checkIfFileExist(completeFileName) {
		fmt.Println("File not found")
	} else {
		data := fileManager.readData(completeFileName)
		dataString := string(data[:])

		if CheckFileValidity(fileName, data){
			fmt.Println(dataString)
		}
	}
}

func CheckFileValidity(id string, data []byte) bool{
	hash := sha256.Sum256(data)
	idFile := NewKademliaIDFromBytes(hash[:IDLength])
	if !idFile.Equals(NewKademliaID(id)){
		fmt.Println(id + " : WARNING ! Modification have been made on this file, it is not valid anymore !")
		return false
	}
	return true
}

func Print(nodesToCheck []NodeToCheck) {
	for i:=0 ; i < len(nodesToCheck) ; i++{
		fmt.Println(nodesToCheck[i].contact.String() + "  alrdyChecked: " + strconv.FormatBool(nodesToCheck[i].alreadyChecked))
	}
}

func (kademlia *Kademlia) GetNetwork() *Network{
	return kademlia.network
}

func (kademlia *Kademlia) GenerateNewFile() string{
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	newData := RandStringRunes(200, letterRunes)
	data := []byte(newData)
	base64Data := base64.StdEncoding.EncodeToString(data[:])

	//Generate a hash for the name of the file
	hash := sha256.Sum256(data)
	idFile := NewKademliaIDFromBytes(hash[:IDLength])
	kademlia.network.FileManager.CheckAndStore(idFile.String(), base64Data)
	return idFile.String()
}

//From internet
func RandStringRunes(n int, letterRunes []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}