package main

import (
	"os"
	"./kademlia"
	"strconv"
	//"time"
	"fmt"
	//"bufio"
	//"strings"
	"time"
	"os/exec"
)

func main() {

	//scanner := bufio.NewScanner(os.Stdin)

	//Take as parameters: id in [0, number of source nodes] number of source nodes, prefix IP, port
	id, _ :=  strconv.Atoi(os.Args[1])
	numberSrcNodes, _ :=  strconv.Atoi(os.Args[2])
	prefixIp := os.Args[3]
	port := os.Args[4]
	var network *kademlia.Network = kademlia.CreateWantedNetwork(id, prefixIp, port)
	kademlia.MakeMoreFriends(network, id, numberSrcNodes, prefixIp, port)
	var kadem *kademlia.Kademlia = &kademlia.Kademlia{}
	kademlia.AssingNetworkKademlia(network, kadem)
	go network.Listen()
	go network.GetMyRoutingTable().StartRoutingTableListener()

	kadem.GetNetwork().GetMyRoutingTable().Print()

	cmd := exec.Command("mkdir", "./kademlia/Files")
	err := cmd.Run()
	if err != nil {
		fmt.Println("$ file list")
	}

	for {
		time.Sleep(10 * time.Second)
		//fmt.Println("supernode")
		//kadem.GetNetwork().GetMyRoutingTable().Print()
		kademlia.ListFiles()
		//fmt.Println("^^^^^^^")
	}

}

