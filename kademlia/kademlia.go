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
	nodesToCheck [bucketSize]NodeToCheck
	mux sync.Mutex
}

//LookupContact is a method of KAdemlia to locate some Node
//PArams target: it is the finded contact
func (kademlia *Kademlia) LookupContact(target *Contact) {
	//Creation of the shared array
	var safeNodesTocheck = SafeNodesCheck{}
	result := false
	var nodesToCheck [bucketSize]NodeToCheck
	nbRunningThreads := 0
	network := Network{}

	//Fulfill this array with the k firt nodes
	var firstContacts []Contact = kademlia.routintTable.FindClosestContacts(target.ID, bucketSize)
	for i := 0; i < len(firstContacts) ; i++ {
		nodesToCheck[i] = NodeToCheck{&firstContacts[i], false}
	}

	//Start sending rpc nodes ->wait method
	for {
		if nbRunningThreads <= alpha{
			nbRunningThreads++
			go safeNodesTocheck.sendFindNode(&nbRunningThreads, &network, &result)
		}

		if(!result){
			break
		}

	}

	//Recursively send back to the next closest one
	//If no respond, don't take into account

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
func(safeNodesCheck *SafeNodesCheck) sendFindNode(nbRunningThreads *int, network *Network, result *bool){
	safeNodesCheck.mux.Lock()
	//find the next one to check
	for i:=0 ; i < bucketSize ; i++{
		if !safeNodesCheck.nodesToCheck[i].alreadyChecked {
			safeNodesCheck.nodesToCheck[i].alreadyChecked = true
			network.SendFindContactMessage(safeNodesCheck.nodesToCheck[i].contact)
			safeNodesCheck.mux.Unlock()
			var newContacts ContactCandidates = network.getResponse()

			//insertion of the new one
			safeNodesCheck.mux.Lock()
			for i:=0 ; i < len(newContacts.contacts) ; i++{
				for j:=0 ; j<bucketSize ; j++{
					if newContacts.contacts[i].Less(safeNodesCheck.nodesToCheck[j].contact) &&
											!newContacts.contacts[i].ID.Equals(safeNodesCheck.nodesToCheck[j].contact.ID) {

						safeNodesCheck.nodesToCheck[j].contact = &(newContacts.contacts[i])
						safeNodesCheck.nodesToCheck[j].alreadyChecked = false
					}
				}
			}
			safeNodesCheck.mux.Unlock()
			*nbRunningThreads--
			*result =  true
		}
	}

	safeNodesCheck.mux.Unlock()
	*nbRunningThreads--
	*result = false
}
