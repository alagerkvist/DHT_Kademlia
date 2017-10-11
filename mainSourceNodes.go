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


const PRINT_FILES = false
const PRINT_ROUTING_TABLE = false
const PRINT_PONG  = true

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
	go kadem.GetNetwork().Listen()
	go kadem.StartRoutingTableListener()
	go kadem.CheckFiles()

	//kadem.GetNetwork().GetMyRoutingTable().Print()

	cmd := exec.Command("mkdir", "./kademlia/Files")
	err := cmd.Run()
	if err != nil {
		fmt.Println("$ file list")
	}

	for {
		time.Sleep(10 * time.Second)
		if(PRINT_ROUTING_TABLE) {
			fmt.Println("supernode")
			kadem.GetNetwork().GetMyRoutingTable().Print()
		}
		if(PRINT_FILES) {
			fmt.Println("^^^ID" + kadem.GetNetwork().GetMyRoutingTable().GetMyContact().ID.String() + "^^^FILES")
			kademlia.ListFiles()
			fmt.Println("*********************")
		}
		//fmt.Println("^^^^^^^")
	}

}

