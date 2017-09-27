package kademlia

import (
	"strconv"
	"fmt"
)

func CreateRandomNetworks(numberNodes int) []Network{
	var newNetworks []Network = make([]Network, numberNodes)
	fmt.Println("creating")
	for i:= 0 ; i < numberNodes ; i++ {
		var newKademliaId *KademliaID = NewRandomKademliaID()
		number := 3000
		var newContact = NewContact(newKademliaId, "127.0.0.1:" + strconv.Itoa(number))
		newNetworks[i].myRoutingTable = NewRoutingTable(newContact)
	}
	return newNetworks
}

func MakeMoreFriends(network *Network, id int, numberSrcNodes int, ipPrefix string, port string){
	for j := 0 ; j < numberSrcNodes ; j++{
		if(j != id){
			newKademliaID := strconv.FormatInt(int64(j), 16) + "00000000000000000000000000000000000000000"
			newContact := NewContact(NewKademliaID(newKademliaID), ipPrefix + strconv.Itoa(20 + j + 1) + ":" + port)
			newContact.CalcDistance(network.myRoutingTable.me.ID)
			network.myRoutingTable.AddContact(newContact)
		}
	}
}


func CreateWantedNetwork(id int, ipPrefix string, port string) *Network{
	var newNetwork *Network = &Network{}

	//creation of the id
	newKademliaID := strconv.FormatInt(int64(id), 16) + "00000000000000000000000000000000000000000"

	var newContact = NewContact(NewKademliaID(newKademliaID), ipPrefix + strconv.Itoa(20 + id + 1) + ":" + port)
	newNetwork.myRoutingTable = NewRoutingTable(newContact)
	return newNetwork
}


func AssingNetworkKademlia(networks []Network, kademlias []Kademlia){
	for i := 0 ; i < len(networks) ; i++{
		kademlias[i].network =  &networks[i]
	}
}