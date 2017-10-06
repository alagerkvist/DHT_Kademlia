package main

import (
	"./kademlia"
	"time"

	"fmt"
)

func main() {


	numberNodes := 16

	//Creation of the Kademlia nodes
	fmt.Println(kademlia.NewRandomKademliaID())

	kademliaNodes := make([]kademlia.Kademlia, numberNodes)
	networks  := kademlia.CreateWantedNetworkPrev(numberNodes)
	//The first node is the new node
	kademlia.MakeMoreFriendsPrev(networks, 5)



	kademlia.AssingNetworkKademliaPrev(networks, kademliaNodes)
	go kademliaNodes[0].GetNetwork().GetMyRoutingTable().StartRoutingTableListener()


	for i:=1 ; i < numberNodes ; i++{
		go kademliaNodes[i].GetNetwork().GetMyRoutingTable().StartRoutingTableListener()
		go kademliaNodes[i].GetNetwork().Listen()
	}

	time.Sleep(2 * time.Second)

	//kademliaNodes[0].GetNetwork().GetMyRoutingTable().Print()
	//kademliaNodes[0].Store("kademlia/routingtable.go")

	//	kademliaNodes[0].LookupContact(kademliaNodes[0].GetNetwork().GetMyRoutingTable().GetMyContact().ID)
	kademliaNodes[0].PrintFile("5cadfe84814b7c5b1f027e0a5dd4d51e89eb6429")
	//time.Sleep(2 * time.Second)
//	kademliaNodes[0].GetNetwork().GetMyRoutingTable().Print()

	//	go kademliaNodes[0].StartRefreshManaging()
//	time.Sleep(250 * time.Second)

}
