package kademlia

import (
	"strconv"
	"fmt"
	"math/rand"
	"encoding/base64"
)

func CreateRandomNetworks(numberNodes int, ip string, port string) *Network{

	var newNetwork Network = Network{}
	var newKademliaId *KademliaID = NewRandomKademliaID()

	var newContact = NewContact(newKademliaId, ip + ":" + port)
	newNetwork.myRoutingTable = NewRoutingTable(newContact)

	return &newNetwork
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

func AddSourceNodes(network *Network, numberSourcesNodes int, ip string, port string){
	ipPrefix := ip[: len(ip) - 1]
	for j := 0 ; j < numberSourcesNodes ; j++{
			newKademliaID := strconv.FormatInt(int64(j), 16) + "00000000000000000000000000000000000000000"
			newContact := NewContact(NewKademliaID(newKademliaID), ipPrefix + strconv.Itoa(20 + j + 1) + ":" + port)
			newContact.CalcDistance(network.myRoutingTable.me.ID)
			network.myRoutingTable.AddContact(newContact)
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


func AssingNetworkKademlia(networks *Network, kademlia *Kademlia){
	kademlia.network =  networks
}
func CreateWantedNetworkPrev(numberNodes int) []Network{
	var newNetworks []Network = make([]Network, numberNodes)
	var ids []string = make([]string, numberNodes)

	//creation of the ids
	for i := 0 ; i < numberNodes ; i++{
		hexad := fmt.Sprintf("%x", i)
		ids[i] = "1111111" + hexad + "00000000000000000000000000000000"
		number := 1234 + i
		var newContact = NewContact(NewKademliaID(ids[i]), "127.0.0.1:" + strconv.Itoa(number))
		newNetworks[i].myRoutingTable = NewRoutingTable(newContact)
		newNetworks[i].fileManager = &FileManager{&base64.Encoding{}, make(map[string]FileInfo)}

}

	return newNetworks

}


func MakeMoreFriendsPrev(nodeToMakeFriends []Network, newFriends int){
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


func AssingNetworkKademliaPrev(networks []Network, kademlias []Kademlia){
	for i := 0 ; i < len(networks) ; i++{
		kademlias[i].network =  &networks[i]
	}
}