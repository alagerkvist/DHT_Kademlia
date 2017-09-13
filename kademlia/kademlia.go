package kademlia

type Kademlia struct {
}

//LookupContact is a method of KAdemlia to locate some Node
//PArams target: it is the finded contact
func (kademlia *Kademlia) LookupContact(target *Contact) {
	//pick alpha nodes from the closest bucket
	//send FIND_NODE to them
	//WARNING bucket and find replies could be augmented with RTT estimates ---> ?

	//Recursively send back to the next closest one
	//If no respond, don't take into account

	// TODO Testing ssh
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
