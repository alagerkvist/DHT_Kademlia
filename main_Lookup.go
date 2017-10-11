package main

import (
	"./kademlia"
	"time"
)

func main() {


	numberNodes := 16

	kademliaNodes := make([]kademlia.Kademlia, numberNodes)
	networks  := kademlia.CreateWantedNetworkPrev(numberNodes)
	//The first node is the new node
	kademlia.MakeMoreFriendsPrev(networks, 5)



	kademlia.AssingNetworkKademliaPrev(networks, kademliaNodes)
	go kademliaNodes[0].StartRoutingTableListener()
	//go kademliaNodes[0].StartRefreshManaging()
	//go kademliaNodes[0].CheckFiles()


	for i:=1 ; i < numberNodes ; i++{
		go kademliaNodes[i].StartRoutingTableListener()
		go kademliaNodes[i].GetNetwork().Listen()

	}
	time.Sleep(2 * time.Second)


	kademliaNodes[0].GetNetwork().GetMyRoutingTable().Print()
	kademliaNodes[0].LookupContact(kademliaNodes[0].GetNetwork().GetMyRoutingTable().GetMyContact().ID)
	kademliaNodes[0].GetNetwork().GetMyRoutingTable().Print()

	//kademliaNodes[0].Store("mainAppl.go")
	//time.Sleep(120 * time.Second)

	//kademliaNodes[0].Store("mainAppl.go")
	//kademliaNodes[0].GetNetwork().GetMyRoutingTable().Print()
	//kademliaNodes[0].Store("kademlia/routingtable.go")

	//	kademliaNodes[0].LookupContact(kademliaNodes[0].GetNetwork().GetMyRoutingTable().GetMyContact().ID)

	//kademliaNodes[0].PrintFile(newFile)
	//kademliaNodes[0].PrintFile(filePrefix + kademliaNodes[0].GenerateNewFile())
	//time.Sleep(2 * time.Second)
//	kademliaNodes[0].GetNetwork().GetMyRoutingTable().Print()

	//	go kademliaNodes[0].StartRefreshManaging()
//	time.Sleep(250 * time.Second)

}
