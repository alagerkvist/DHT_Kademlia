package kademlia

import (
	"sync"
	"fmt"
	"strconv"
	"crypto/sha256"
	"encoding/base64"
)

const alpha = 3

type Kademlia struct {
	network *Network
}

type NodeToCheck struct{
	contact	*Contact
	alreadyChecked	bool 	//usefull if network < k
}

type SafeNodesCheck struct {
	nodesToCheck []NodeToCheck
	mux sync.Mutex
}

type Request struct{
	contact *Contact
	endWork bool
}

type Response struct{
	newContacts *ContactCandidates
	data *string
}

//LookupContact is a method of Kademlia to locate some Node
//Params target: it is the finded contact
func (kademlia *Kademlia) LookupContact(targetId *KademliaID) []NodeToCheck{
	//Creation of the shared array
	var safeNodesToCheck = SafeNodesCheck{nodesToCheck: make([]NodeToCheck, 0, bucketSize)}
	var noMoreToCheck bool
	nbRunningThreads := 0
	var wg sync.WaitGroup
	var counterThreadsDone = 0

	//Fulfill this array with at most the k nodes from buckets
	responseChannel := make(chan []Contact)
	kademlia.network.myRoutingTable.createTask(lookUpContact, responseChannel, &Contact{targetId, "", nil})
	firstContacts := <- responseChannel

	//fmt.Println("\n--- Contacts to join: ")
	for i:=0 ; i<len(firstContacts) ; i++ {
		//fmt.Println(firstContacts[i].String())
		safeNodesToCheck.nodesToCheck = append(safeNodesToCheck.nodesToCheck, NodeToCheck{&firstContacts[i], false})
	}
	//fmt.Println("---")

	//Start sending rpc nodes -> wait method
	for {
		if nbRunningThreads < alpha {
			wg.Add(1)
			counterThreadsDone++
			nbRunningThreads++
			go safeNodesToCheck.sendFindNode(&nbRunningThreads, kademlia.network, &noMoreToCheck, &wg, targetId)
		}

		//Stop the infinite loop when get all the information needed
		if(noMoreToCheck){
			wg.Wait() //Wait for all threads to finish
			if(noMoreToCheck){
				break
			}
		}
	}

	return safeNodesToCheck.nodesToCheck
}

/*
*	Send RPC_NODE RPC and add new contacts to the array, check if need to end the loop.
 */
func(safeNodesCheck *SafeNodesCheck) sendFindNode(nbRunningThreads *int, network *Network, noMoreToCheck *bool, wg *sync.WaitGroup, targetID *KademliaID){
	defer wg.Done()
	safeNodesCheck.mux.Lock()
	//find the next one to check
	for i:=0 ; i < len(safeNodesCheck.nodesToCheck) ; i++{
		if !safeNodesCheck.nodesToCheck[i].alreadyChecked && !safeNodesCheck.nodesToCheck[i].contact.ID.Equals(network.myRoutingTable.me.ID){
			//The node will not be taken into account by the other threads.
			safeNodesCheck.nodesToCheck[i].alreadyChecked = true
			safeNodesCheck.mux.Unlock()

			//fmt.Println("\n*** Ask to: " + safeNodesCheck.nodesToCheck[i].contact.String())
			//fmt.Println("Target: "+ targetID.String() +"\n***\n")

			newContacts := network.SendFindContactMessage(safeNodesCheck.nodesToCheck[i].contact, targetID)

			//insertion of the new one
			safeNodesCheck.mux.Lock()
			for i:=0 ; i<newContacts.Len() ; i++{
				//fmt.Println("Contact to add: " + newContacts.contacts[i].String())

				for j:=0 ; j<bucketSize ; j++{
					if !newContacts.contacts[i].ID.Equals(network.myRoutingTable.me.ID){
						//case of less than k values
						if j >= len(safeNodesCheck.nodesToCheck) {
							safeNodesCheck.nodesToCheck = append(safeNodesCheck.nodesToCheck, NodeToCheck{&newContacts.contacts[i], false})
							break;
						} else if newContacts.contacts[i].ID.Equals(safeNodesCheck.nodesToCheck[j].contact.ID){
							break;
						} else if  newContacts.contacts[i].Less(safeNodesCheck.nodesToCheck[j].contact){
							//Shift value of the array and insert the new one
							copy(safeNodesCheck.nodesToCheck[j+1:], safeNodesCheck.nodesToCheck[j: len(safeNodesCheck.nodesToCheck) - 1])
							safeNodesCheck.nodesToCheck[j].contact = &(newContacts.contacts[i])
							safeNodesCheck.nodesToCheck[j].alreadyChecked = false
							break;
						}
					}
				}
			}
			safeNodesCheck.mux.Unlock()
			*nbRunningThreads--
			*noMoreToCheck = false


			/*fmt.Println("\n%%%% After adding: ")
			safeNodesCheck.Print()
			fmt.Println("\n%%%")
			*/
			return
		}
	}

	safeNodesCheck.mux.Unlock()
	*nbRunningThreads--
	*noMoreToCheck = true

	fmt.Println("\n%%%% Nothing to add \n%%%")

	return
}



//LookupContact is a method of KAdemlia to locate some Data
//PArams hash: it is the finded data with the 160 bits hash
func (kademlia *Kademlia) LookupData(fileName string) {
	targetID := NewKademliaID(fileName)
	channelToSendRequest := make(chan Request, alpha)
	channelToReceive := make(chan Response)
	nodesToCheck := make([]NodeToCheck, 0, bucketSize)

	//Fulfill this array with at most the k nodes from buckets
	responseChannel := make(chan []Contact)
	kademlia.network.myRoutingTable.createTask(lookUpContact, responseChannel, &Contact{targetID, "", nil})
	firstContacts := <- responseChannel

	for i:=0 ; i<len(firstContacts) ; i++ {
		nodesToCheck = append(nodesToCheck, NodeToCheck{&firstContacts[i], false})
	}

	for i := 0 ; i < alpha ; i++ {
		go kademlia.network.workerFindData(channelToSendRequest, *targetID, channelToReceive)
		channelToSendRequest <- Request{nodesToCheck[i].contact, false}
		nodesToCheck[i].alreadyChecked = true
	}

	for {
		newResponse := <- channelToReceive
		//Receive data
		if(newResponse.newContacts == nil){
			fmt.Println("response: " + *newResponse.data)
			kademlia.network.fileManager.checkAndStore(fileName, *newResponse.data)
			sendEndWork(channelToSendRequest, alpha)
			break
		}

		newContacts := newResponse.newContacts

		for i:=0 ; i<newContacts.Len() ; i++{
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

		//No nodes with this data
		nextContactToCheck :=  kademlia.network.getNextContactToAsk(nodesToCheck)
		if nextContactToCheck == nil{
			fmt.Println("Impossible to find the file in the network")
			Print(nodesToCheck)
			sendEndWork(channelToSendRequest, alpha)
			break
		}
		channelToSendRequest <- Request{kademlia.network.getNextContactToAsk(nodesToCheck), false}
	}


}


func sendEndWork(channelToSendRequest chan Request, nb int){
	for i := 0 ; i < alpha ; i++{
		channelToSendRequest <- Request{nil, true}
	}
}


func(network *Network) workerFindData(requestsChannel chan Request, targetId KademliaID, responseChannel chan Response) {

	for {
		request := <- requestsChannel
		if(request.endWork){
			break
		}
		responseChannel <- network.SendFindDataValue(targetId, request.contact)
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
	fileManager := kademlia.network.fileManager

	if !fileManager.checkIfFileExist(fileName){
		fmt.Println("File not found")
	}else {
		data := fileManager.readData(fileName)
		base64Data := base64.StdEncoding.EncodeToString(data[:])

		//Generate a hash for the name of the file
		hash := sha256.Sum256(data)
		idFile := NewKademliaIDFromBytes(hash[:IDLength])

		fileManager.checkAndStore(idFile.String(), base64Data)

		contactToSend := kademlia.LookupContact(idFile)
		fmt.Println(contactToSend)
		kademlia.network.SendStoreMessage(idFile.String(), base64Data, contactToSend)
	}
}

func Print(nodesToCheck []NodeToCheck) {
	for i:=0 ; i < len(nodesToCheck) ; i++{
		fmt.Println(nodesToCheck[i].contact.String() + "  alrdyChecked: " + strconv.FormatBool(nodesToCheck[i].alreadyChecked))
	}
}

func (kademlia *Kademlia) GetNetwork() *Network{
	return kademlia.network
}


