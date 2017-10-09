package kademlia

import (
	"strconv"
	"fmt"
	"math/rand"
	"strings"
)

/**
*	CreateRandomNetworks
* numberNodes: the number of big nodes
* ip: the ip
* port: the port of the nodes
* Create a node
 */
func CreateRandomNetworks(numberNodes int, ip string, port string) *Network{

	var newNetwork Network = Network{}
	var newKademliaId *KademliaID = NewRandomKademliaID()
	fmt.Println(newKademliaId)
	var newContact = NewContact(newKademliaId, ip + ":" + port)
	newNetwork.myRoutingTable = NewRoutingTable(newContact)
	newNetwork.FileManager = &FileManager{}
	newNetwork.FileManager.filesStored = make(map[string]*FileInfo)

	return &newNetwork
}

/**
*	MakeMoreFriends
* network
* id: the number of important nodes
* numberSrcNodes: the ip
* ipPrefix: the prefix of the ip
* port: the port
* Add some nodes into the routing table
 */
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

/**
*	AddSourceNodes
* network
* ip: ip of the node
* port: port of this node
* Add the main nodes into the routing table
 */
func AddSourceNodes(network *Network, numberSourcesNodes int, ip string, port string){
	ipPrefixAux := strings.Split(ip,".")
	ipPrefix := ipPrefixAux[0] +"."+ipPrefixAux[1]+"."+ipPrefixAux[2]+"."
	for j := 0 ; j < numberSourcesNodes ; j++{
			newKademliaID := strconv.FormatInt(int64(j), 16) + "00000000000000000000000000000000000000000"
			newContact := NewContact(NewKademliaID(newKademliaID), ipPrefix + strconv.Itoa(20 + j + 1) + ":" + port)
			newContact.CalcDistance(network.myRoutingTable.me.ID)
			network.myRoutingTable.AddContact(newContact)
	}
}


/**
*	CreateWantedNetwork
* network
* id: the number of important nodes
* numberSrcNodes: the ip
* ipPrefix: the prefix of the ip
* port: the port
* Add the main nodes into the routing table
 */
func CreateWantedNetwork(id int, ipPrefix string, port string) *Network{
	var newNetwork *Network = &Network{}

	//creation of the id
	newKademliaID := strconv.FormatInt(int64(id), 16) + "00000000000000000000000000000000000000000"

	var newContact = NewContact(NewKademliaID(newKademliaID), ipPrefix + strconv.Itoa(20 + id + 1) + ":" + port)
	newNetwork.myRoutingTable = NewRoutingTable(newContact)
	newNetwork.FileManager = &FileManager{make(map[string]*FileInfo)}
	return newNetwork
}


/**
*	AssingNetworkKademlia
* network
* id: the number of important nodes
* numberSrcNodes: the ip
* ipPrefix: the prefix of the ip
* port: the port
* Add the main nodes into the routing table
 */
func AssingNetworkKademlia(networks *Network, kademlia *Kademlia){
	kademlia.network =  networks
}


/**
*	CreateWantedNetworkPrev
* network
* id: the number of important nodes
* numberSrcNodes: the ip
* ipPrefix: the prefix of the ip
* port: the port
* Add the main nodes into the routing table
 */
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
		newNetworks[i].FileManager = &FileManager{make(map[string]*FileInfo)}

}
	return newNetworks
}


/**
*	MakeMoreFriendsPrev
* network
* id: the number of important nodes
* numberSrcNodes: the ip
* ipPrefix: the prefix of the ip
* port: the port
* Add the main nodes into the routing table
 */
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


/**
*	AssingNetworkKademliaPrev
*   assign network to kademlia structures
 */
func AssingNetworkKademliaPrev(networks []Network, kademlias []Kademlia){
	for i := 0 ; i < len(networks) ; i++{
		kademlias[i].network =  &networks[i]
	}
}
