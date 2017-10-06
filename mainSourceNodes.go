package main

import (
	"os"
	"./kademlia"
	"strconv"
	//"time"
	"fmt"
	//"bufio"
	"strings"
	"time"
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

	for {
		time.Sleep(10 * time.Second)
		fmt.Println("supernode")
		kadem.GetNetwork().GetMyRoutingTable().Print()
		fmt.Println("^^^^^^^")
	}
/*
	for scanner.Scan() {
		processText_SN(scanner.Text(), kadem)

	}*/
	/*
	fmt.Println(len(os.Args))
	if len(os.Args) <= 5 {

		fmt.Println("supernode leng")
	}else{
		printHelp_SN()

	}
*/

}


func processText_SN(text string, kadem *kademlia.Kademlia){
	//fmt.Println(text)
	//fmt.Println("hello")
	var words []string = strings.Split(text," ");

	var command = words[0]
	switch command {
	case "ping":
		go processCommandPing_SN(words, kadem)
		break
	case "info":
		go processCommandInfo_SN(words)
		break
	case "routingTable":
		go processCommandRoutingTable_SN(words, kadem)
		break
	case "lookup":
		go processCommandLookup_SN(words, kadem)
		break
	case "help":
		break
	default:
		fmt.Println("command not found")
	}
	printHelp_SN()
}

func processCommandLookup_SN (words []string, kadem *kademlia.Kademlia){
	fmt.Println("processCommandLookup_SN", kadem.GetNetwork().GetMyRoutingTable().GetMyContact().ID)
	kadem.LookupContact(kadem.GetNetwork().GetMyRoutingTable().GetMyContact().ID)
}


func processCommandPing_SN(words []string, kadem *kademlia.Kademlia){
	if  len(words) != 3 ||
		(words[1] != "--nodeID" && words[1] != "--nodeIP") {
		fmt.Println("error PING")
		return
	}


	newKademliaID := strconv.FormatInt(int64(0), 16) + "00000000000000000000000000000000000000000"
	newContact := kademlia.NewContact(kademlia.NewKademliaID(newKademliaID), "10.5.0." + strconv.Itoa(21) + ":" + strconv.Itoa(8080))
	kadem.GetNetwork().SendPingMessage(&newContact)

	fmt.Println("sending ping to ", words[2])
	fmt.Println("***********")
}

func processCommandInfo_SN(words []string)  {
	if  len(words) != 1 || words[0] != "info" {
		fmt.Println("error INFO")
		return
	}
	fmt.Println("my info")
	fmt.Println("***********")
}

func processCommandRoutingTable_SN(words []string, kadem *kademlia.Kademlia){
	if  len(words) != 1 || words[0] != "routingTable" {
		fmt.Println("error routingTable")
		return
	}

	fmt.Println("my routing table")
	kadem.GetNetwork().GetMyRoutingTable().Print()
	fmt.Println("***********")
}

func printHelp_SN(){
	fmt.Println("*****************************************")
	fmt.Println("This is my help: (write option + enter)")
	fmt.Println("")
	fmt.Println("PING make a ping to the selected node")
	fmt.Println("$ ping --nodeID KademliaID")
	fmt.Println("$ ping --nodeIP KademliaID")
	fmt.Println("")
	fmt.Println("Return info about the current node")
	fmt.Println("$ info")
	fmt.Println("")
	fmt.Println("Return the routing table")
	fmt.Println("$ routingTable")
	fmt.Println("*****************************************")
}
