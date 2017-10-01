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
func (kademlia *Kademlia) LookupData(hash string) {

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




func (safeNodeToCheck *SafeNodesCheck) Print() {
	for i:=0 ; i < len(safeNodeToCheck.nodesToCheck) ; i++{
		fmt.Println(safeNodeToCheck.nodesToCheck[i].contact.String() + "  alrdyChecked: " + strconv.FormatBool(safeNodeToCheck.nodesToCheck[i].alreadyChecked))
	}
}

func (kademlia *Kademlia) GetNetwork() *Network{
	return kademlia.network
}


