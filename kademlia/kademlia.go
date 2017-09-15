package kademlia

import "sync"

const alpha = 3

type Kademlia struct {
	routintTable RoutingTable
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
func (kademlia *Kademlia) LookupContact(target *Contact) {
	//Creation of the shared array
	var safeNodesToCheck = SafeNodesCheck{nodesToCheck: make([]NodeToCheck, 0, bucketSize)}
	var noMoreToCheck bool
	nbRunningThreads := 0
	network := Network{}
	var wg sync.WaitGroup
	var counterThreadsDone = 0

	//Fulfill this array with at most the k nodes from buckets
	var firstContacts []Contact = kademlia.routintTable.FindClosestContacts(target.ID, bucketSize)
	for i:=0 ; i<len(firstContacts) ; i++ {
		safeNodesToCheck.nodesToCheck = append(safeNodesToCheck.nodesToCheck, NodeToCheck{&firstContacts[i], false})
	}

	//Start sending rpc nodes ->wait method
	for {
		if nbRunningThreads <= alpha {
			wg.Add(1)
			counterThreadsDone++
			nbRunningThreads++
			go safeNodesToCheck.sendFindNode(&nbRunningThreads, &network, &noMoreToCheck, &wg)
		}

		//Stop the infinite loop when get all the information needed
		if(noMoreToCheck){
			wg.Wait() //Wait for all threads to finish
			if(noMoreToCheck){
				break
			}
		}
	}
}

//LookupContact is a method of KAdemlia to locate some Data
//PArams hash: it is the finded data with the 160 bits hash
func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

//Store is the method of KAdemlia to Store data
// Params: data array of Bytes.
func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}


/*
*	Send RPC_NODE RPC and add new contacts to the array, check if need to end the loop.
 */
func(safeNodesCheck *SafeNodesCheck) sendFindNode(nbRunningThreads *int, network *Network, noMoreToCheck *bool, wg *sync.WaitGroup){
	defer wg.Done()
	safeNodesCheck.mux.Lock()
	//find the next one to check
	for i:=0 ; i < len(safeNodesCheck.nodesToCheck) ; i++{
		if !safeNodesCheck.nodesToCheck[i].alreadyChecked {
			//The node will not be taken into account by the other threads.
			safeNodesCheck.nodesToCheck[i].alreadyChecked = true
			safeNodesCheck.mux.Unlock()
			network.SendFindContactMessage(safeNodesCheck.nodesToCheck[i].contact)
			var newContacts ContactCandidates = network.getResponse()

			//insertion of the new one
			safeNodesCheck.mux.Lock()
			for i:=0 ; i<len(newContacts.contacts) ; i++{
				for j:=0 ; j<bucketSize ; j++{
					//case of less than k values
					if j >= len(safeNodesCheck.nodesToCheck) {
						safeNodesCheck.nodesToCheck = append(safeNodesCheck.nodesToCheck, NodeToCheck{&newContacts.contacts[i], false})

					} else if  newContacts.contacts[i].Less(safeNodesCheck.nodesToCheck[j].contact) &&
							!newContacts.contacts[i].ID.Equals(safeNodesCheck.nodesToCheck[j].contact.ID){
						//Shift value of the array and insert the new one
						copy(safeNodesCheck.nodesToCheck[j+1:], safeNodesCheck.nodesToCheck[j: len(safeNodesCheck.nodesToCheck) - 1])
						safeNodesCheck.nodesToCheck[j].contact = &(newContacts.contacts[i])
						safeNodesCheck.nodesToCheck[j].alreadyChecked = false
					}
				}
			}
			safeNodesCheck.mux.Unlock()
			*nbRunningThreads--
			*noMoreToCheck = false
			return
		}
	}

	safeNodesCheck.mux.Unlock()
	*nbRunningThreads--
	*noMoreToCheck = true
	return
}
