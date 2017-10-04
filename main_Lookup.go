package main

import (
	"./kademlia"
	"time"
)

func main() {


	numberNodes := 15

	//Creation of the Kademlia nodes
	kademliaNodes := make([]kademlia.Kademlia, numberNodes)
	networks  := kademlia.CreateWantedNetworkPrev(numberNodes)
	//The first node is the new node
	kademlia.MakeMoreFriendsPrev(networks, 10)



	kademlia.AssingNetworkKademliaPrev(networks, kademliaNodes)
	go kademliaNodes[0].GetNetwork().GetMyRoutingTable().StartRoutingTableListener()


	for i:=1 ; i < numberNodes ; i++{
		go kademliaNodes[i].GetNetwork().GetMyRoutingTable().StartRoutingTableListener()
		go kademliaNodes[i].GetNetwork().Listen()
	}
	time.Sleep(2 * time.Second)

	kademliaNodes[0].LookupContact(kademliaNodes[0].GetNetwork().GetMyRoutingTable().GetMyContact().ID)
	//kademliaNodes[0].LookupData("98c52cbb1057afa0af21d602a0c5ccde4a762d0a")

}
