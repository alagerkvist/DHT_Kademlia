package main

import (
	"os"
	"./kademlia"
	"strconv"
	"time"
	"fmt"
)

func main() {

	//Take as parameters: id in [0, number of source nodes] number of source nodes, prefix IP, port
	id, _ :=  strconv.Atoi(os.Args[1])
	numberSrcNodes, _ :=  strconv.Atoi(os.Args[2])
	prefixIp := os.Args[3]
	port := os.Args[4]
	var network *kademlia.Network = kademlia.CreateWantedNetwork(id, prefixIp, port)
	kademlia.MakeMoreFriends(network, id, numberSrcNodes, prefixIp, port)
	go network.Listen()
	for {
		time.Sleep(20 * time.Second)
		fmt.Println("supernode")
	}
}


