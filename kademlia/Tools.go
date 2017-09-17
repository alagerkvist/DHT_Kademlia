package kademlia

import (
 "math/rand"
)

func CreateRandomNetworks(numberNodes int) []Network{
	var newNetworks []Network = make([]Network, numberNodes)

	for i:= 0 ; i < numberNodes ; i++{
		var newKademliaId *KademliaID = NewRandomKademliaID()
		number := 8000 + i
		*newNetworks[i].myContact = NewContact(NewRandomKademliaID(), "localhost:" + string(number))
		newNetworks[i].myRoutingTable = NewRoutingTable(NewContact(newKademliaId, "localhost:" + string(8000 + i)))
	}
	return newNetworks
}

func MakeMoreFriends(nodeToMakeFriends []Network, newFriends int){
	for j := 0 ; j < len(nodeToMakeFriends) ; j++{
		for i := 0 ; i < newFriends ; i++{
			random := rand.Intn(len(nodeToMakeFriends))
			nodeToMakeFriends[j].myRoutingTable.AddContact(*nodeToMakeFriends[random].myContact)
		}
	}
}