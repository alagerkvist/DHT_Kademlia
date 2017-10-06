package main

import (
	"strconv"
	"os"
	"time"
	"fmt"
	"./kademlia"
	"bufio"
)

func main() {

	//Take as parameters: id in [0, number of source nodes] number of source nodes, prefix IP, port
	numberSrcNodes, _ :=  strconv.Atoi(os.Args[1])
	ip := getMyIp()
	port := "8080"

	var network *kademlia.Network = kademlia.CreateRandomNetworks(numberSrcNodes, ip, port)
	kademlia.AddSourceNodes(network, numberSrcNodes, ip, port)
	fmt.Println(network.GetMyRoutingTable().GetMyContact())
	fmt.Println(network.GetMyRoutingTable())
	var kadem *kademlia.Kademlia = &kademlia.Kademlia{}
	kademlia.AssingNetworkKademlia(network, kadem)

	go network.Listen()
	go network.GetMyRoutingTable().StartRoutingTableListener()
	//go kadem.StartRefreshManaging()

}


