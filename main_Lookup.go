package main

import (
	"./kademlia"
	"time"
	"fmt"
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
	fmt.Println("ok")
	time.Sleep(2 * time.Second)

	kademliaNodes[0].LookupContact(kademliaNodes[0].GetNetwork().GetMyRoutingTable().GetMyContact())
}
