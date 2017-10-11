package kademlia

import (
	"fmt"
	"strconv"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"time"
)

//Alpha = the number of workers while a lookup
const alpha = 3

//Main structure of a node.
type Kademlia struct {
	network *Network
}

/** NodeToCheck
*	contact: contact to check
*   alreadyChecked: already sent it a lookup
 */
type NodeToCheck struct{
	contact	*Contact
	alreadyChecked	bool 	//usefull if network < k
}

/** Request: used to communicate with the Lookup workers
* contact: the contact to request
* endWork: prevent the worker that it can stop to work
*/
type Request struct{
	contact *Contact
	endWork bool
}

/** Response: use by the Lookup workers to answer
*  newContacts: the new contacts retrieved from the asked contact
*  data: the data of the file (case of lookupValue
* contactedContact *Contact: the contact whose the request has been performed
* error: flag to signal error during the communication
*/
type Response struct{
	newContacts *ContactCandidates
	data *string
	contactedContact *Contact
	error bool
}

/** LookUpContact
* PARAM: kademlia
*		 targetID: the target id
* OUTPUT: the k closest nodes from the target id
*/
func (kademlia *Kademlia) LookupContact(targetId *KademliaID) []NodeToCheck{
	return kademlia.Lookup(targetId, true)
}

/** LookUpData
* PARAM: fileManager
*		 fileName: the file to retrieve from the network
* Retrieve the file from the network if it is still there
*/
func (kademlia *Kademlia) LookupData(fileName string) {
	targetID := NewKademliaID(fileName)
	kademlia.Lookup(targetID, false)

}

/** LookUp
* PARAM: targetId: the target
		 isForNode: need to stop when receive file or not
* OUTPUT: the k-closest nodes of the target
*/
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


/** SendEndWork
* PARAM: channelToSendRequest: the channel to send request
		 nb: the of thread to finish
* Say to the worker that they can't stop waiting for request
*/
func sendEndWork(channelToSendRequest chan Request, nb int){
	for i := 0 ; i < alpha ; i++{
		channelToSendRequest <- Request{nil, true}
	}
}

/** workerFindData
* PARAM: network
*		 targetId: the target ID
		 responseChannel: channel to put answers
		 isForNode: lookupValue or data
* Perform tasks from the request channel
*/
func(network *Network) workerFindData(requestsChannel chan Request, targetId KademliaID, responseChannel chan Response, isForNode bool) {

	for {
		request := <- requestsChannel
		if(request.endWork){
			break
		}
		//fmt.Print("request: ")
		//fmt.Println(request.contact)
		if !isForNode {
			responseChannel <- network.SendFindDataValue(targetId, request.contact)
		} else {
			responseChannel <- network.SendFindContactMessage(&targetId, request.contact)
		}
	}
}


/** getNextContactToAsk
* PARAM: network
		 nodesToCheck: the nodes to check
* Return the next contact to check
*/
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

	if !fileManager.checkIfFileExist(fileName){
		fmt.Println("File not found")
	}else {
		data := fileManager.readData(fileName)
		base64Data := base64.StdEncoding.EncodeToString(data[:])

		//Generate a hash for the name of the file
		hash := sha256.Sum256(data)
		idFile := NewKademliaIDFromBytes(hash[:IDLength])

		fileManager.CheckAndStore(idFile.String(), base64Data)
		fileInfo := fileManager.filesStored[idFile.String()]
		fmt.Println(fileInfo)
		fileInfo.originalStore = true
		fileInfo.immutable = true
		//fmt.Println(fileManager.filesStored[idFile.String()])

		contactToSend := kademlia.LookupContact(idFile)
		fmt.Println("File will be send to these contacts:")
		Print(contactToSend)
		kademlia.network.SendStoreMessage(idFile.String(), base64Data, contactToSend)
	}
}

/** PrintFile
* PARAM: network
		 fileName: the name of the file
* Print the data of the file
*/
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


/** PrintFile
* PARAM: network
		 fileName: the name of the file
* Print the data of the file
*/
func CheckFileValidity(id string, data []byte) bool{
	hash := sha256.Sum256(data)
	idFile := NewKademliaIDFromBytes(hash[:IDLength])
	if !idFile.Equals(NewKademliaID(id)){
		fmt.Println(id + " : WARNING ! Modification have been made on this file, it is not valid anymore !")
		return false
	}
	return true
}


/** Print the nodes to check
*/
func Print(nodesToCheck []NodeToCheck) {
	for i:=0 ; i < len(nodesToCheck) ; i++{
		fmt.Println(nodesToCheck[i].contact.String() + "  alrdyChecked: " + strconv.FormatBool(nodesToCheck[i].alreadyChecked))
	}
}


/** return the network
*/
func (kademlia *Kademlia) GetNetwork() *Network{
	return kademlia.network
}

/** GenerateNewFile
* PARAM: kademlia
* Create a random file
*/
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


/** CheckFiles
*   PARAM: kademlia
*	Check the files for refreshing and expiration time
*/
func (kademlia *Kademlia) CheckFiles(){

	for{
		time.Sleep(1 * time.Minute)
		fmt.Println("Check for files")
		for k, file := range kademlia.network.FileManager.filesStored {

			//Refreshing each file that has not been refresh from one hour
			if time.Since(file.lastTimeRefreshed).Hours() >= 1 {
				kademlia.Store(k)

				//Refreshing files owned, each 24h
			} else if file.originalStore && time.Since(file.lastOriginalRefreshedStored).Hours() >= 24{
				file.lastOriginalRefreshedStored = time.Now().Local()
				kademlia.Store(k)

				//Delete expirated files
			} else if !file.immutable && time.Since(file.initialStore).Hours() >= file.expirationTime{
				kademlia.network.FileManager.RemoveFile(k)
			}
		}
	}
}

func (kademlia *Kademlia) StartRoutingTableListener() {
	kademlia.network.myRoutingTable.channelTasks = make(chan Task, nb_task_managed)
	go kademlia.runWorker(kademlia.network.myRoutingTable.channelTasks)
}