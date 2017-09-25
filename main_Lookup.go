package main
/*
import (
	"./kademlia"
	"time"
)

func main() {
	numberNodes := 15

	//Creation of the Kademlia nodes
	kademliaNodes := make([]kademlia.Kademlia, numberNodes)
	networks  := kademlia.CreateWantedNetwork(numberNodes)
	//The first node is the new node
	kademlia.MakeMoreFriends(networks, 10)



	kademlia.AssingNetworkKademlia(networks, kademliaNodes)

	for i:=1 ; i < numberNodes ; i++{
		go kademliaNodes[i].GetNetwork().Listen()
	}
	time.Sleep(2 * time.Second)

	kademliaNodes[0].LookupContact(kademliaNodes[0].GetNetwork().GetMyRoutingTable().GetMyContact())
}
*/