package kademlia

import (
 "math/rand"
	"strconv"
	"fmt"
)

func CreateRandomNetworks(numberNodes int) []Network{
	var newNetworks []Network = make([]Network, numberNodes)
	fmt.Println("creating")
	for i:= 0 ; i < numberNodes ; i++ {
		var newKademliaId *KademliaID = NewRandomKademliaID()
		number := 1234 + i
		fmt.Println(number)
		var newContact = NewContact(newKademliaId, "127.0.0.1:" + strconv.Itoa(number))
		newNetworks[i].myRoutingTable = NewRoutingTable(newContact)
	}
	return newNetworks
}

func MakeMoreFriends(nodeToMakeFriends []Network, newFriends int){
	for j := 0 ; j < len(nodeToMakeFriends) ; j++{

		for i := 0 ; i < newFriends ; i++{
			random := j
			for random == j{
				random = rand.Intn(len(nodeToMakeFriends))
			}
			nodeToMakeFriends[j].myRoutingTable.AddContact(nodeToMakeFriends[random].myRoutingTable.me)
		}
	}
}


func CreateWantedNetwork(numberNodes int) []Network{
	var newNetworks []Network = make([]Network, numberNodes)
	var ids []string = make([]string, numberNodes)

	//creation of the ids
	for i := 0 ; i < numberNodes ; i++{
		hexad := fmt.Sprintf("%x", i)
		ids[i] = "1111111" + hexad + "00000000000000000000000000000000"
		number := 1234 + i
		var newContact = NewContact(NewKademliaID(ids[i]), "127.0.0.1:" + strconv.Itoa(number))
		newNetworks[i].myRoutingTable = NewRoutingTable(newContact)

	}

	return newNetworks

}


func AssingNetworkKademlia(networks []Network, kademlias []Kademlia){
	for i := 0 ; i < len(networks) ; i++{
		kademlias[i].network =  &networks[i]
	}
}