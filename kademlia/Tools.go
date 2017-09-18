package kademlia

import (
 "math/rand"
	"strconv"
)

func CreateRandomNetworks(numberNodes int) []Network{
	var newNetworks []Network = make([]Network, numberNodes)

	for i:= 0 ; i < numberNodes ; i++ {
		var newKademliaId *KademliaID = NewRandomKademliaID()
		number := 8000 + i
		var newContact = NewContact(NewRandomKademliaID(), "localhost:"+strconv.Itoa(number))
		newNetworks[i].myContact = &newContact
		newNetworks[i].myRoutingTable = NewRoutingTable(NewContact(newKademliaId, "localhost:"+string(8000+i)))
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