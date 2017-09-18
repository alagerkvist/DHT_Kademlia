package kademlia

import (
 "math/rand"
	"strconv"
)

func CreateRandomNetworks(numberNodes int) []Network{
	var newNetworks []Network = make([]Network, numberNodes)

	for i:= 0 ; i < numberNodes ; i++ {
		var newKademliaId *KademliaID = NewRandomKademliaID()
		number := 1234 + i
		var newContact = NewContact(NewRandomKademliaID(), "127.0.0.1:"+strconv.Itoa(number))
		newNetworks[i].myContact = &newContact
		newNetworks[i].myRoutingTable = NewRoutingTable(NewContact(newKademliaId, "127.0.0.1:"+strconv.Itoa(number)))
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