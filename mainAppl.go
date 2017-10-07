package main

import (
	"./kademlia"
	"fmt"
	"bufio"
	"os"
	"strings"
	"net"

	"strconv"
	//"mime"
	"os/exec"
	//"log"
	"time"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	numberSrcNodes, _ :=  strconv.Atoi(os.Args[1])
	ip := getMyIp()
	port := "8080"

	var network *kademlia.Network = kademlia.CreateRandomNetworks(numberSrcNodes, ip, port)
	kademlia.AddSourceNodes(network, numberSrcNodes, ip, port)
	fmt.Println(network.GetMyRoutingTable().GetMyContact())
	//fmt.Println(network.GetMyRoutingTable())
	var kadem *kademlia.Kademlia = &kademlia.Kademlia{}
	kademlia.AssingNetworkKademlia(network, kadem)

	go network.Listen()
	go network.GetMyRoutingTable().StartRoutingTableListener()

	//kadem.GetNetwork().GetMyRoutingTable().Print()
	go kadem.LookupContact(kadem.GetNetwork().GetMyRoutingTable().GetMyContact().ID)

	cmd := exec.Command("mkdir", "./kademlia/Files")
	err := cmd.Run()
	if err != nil {
		fmt.Println("$ file list")
	}


	printHelp()

	go printTimmerMyListOFFiles(kadem.GetNetwork().GetMyRoutingTable().GetMyContact().ID.String())
	for scanner.Scan() {
		processText(scanner.Text(), kadem)
	}
}


func printTimmerMyListOFFiles(ID string){
	for {
		time.Sleep(10 * time.Second)
		//fmt.Println("supernode")
		//kadem.GetNetwork().GetMyRoutingTable().Print()
		fmt.Println("^^^ID"+ID+"^^^FILES")
		kademlia.ListFiles()
		fmt.Println("*********************")
		//fmt.Println("^^^^^^^")
	}
}
func getMyIp() string{
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//os.Stdout.WriteString(ipnet.IP.String() + "\n")
				return ipnet.IP.String()
			}
		}
	}
	return ""
}


func processText(text string, kadem *kademlia.Kademlia){
	//fmt.Println(text)
	//fmt.Println("hello")
	var words []string = strings.Split(text," ");

	var command = words[0]
	switch command {
	case "ping":
		processCommandPing(words, kadem)
		break
	case "info":
		processCommandInfo(words)
		break
	case "routingTable":
		processCommandRoutingTable(words, kadem)
		break
	case "lookup":
		processCommandLookup(words, kadem)
		break
	case "file":
		processCommandFile(words,kadem)
	case "help":
		break
	default:
		fmt.Println("command not found")
	}
	//printHelp()
}

func processCommandLookup (words []string, kadem *kademlia.Kademlia){
	fmt.Println("processCommandLookup", kadem.GetNetwork().GetMyRoutingTable().GetMyContact().ID)
	kadem.LookupContact(kadem.GetNetwork().GetMyRoutingTable().GetMyContact().ID)
}


func processCommandPing(words []string, kadem *kademlia.Kademlia){
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

func takeName(words []string) string{
	return strings.Split(words[2],"=")[1]
}
func processCommandFile(words []string, kadem *kademlia.Kademlia){

	var command = words[1]
	switch command {
	case "new":
		//fmt.Println("$ file new")
		kadem.GenerateNewFile()
		break
	case "save":
		//fmt.Println("$ file save --name=NAME")
		kadem.Store(takeName(words))
		break
	case "ping":
		//fmt.Println("$ file ping --name=NAME")
		break
	case "take":
		//fmt.Println("$ file find --name=NAME")
		kadem.LookupData(takeName(words))
		break
	case "rm":
		//fmt.Println("$ file rm --name=NAME")
		break
	case "print":
		//fmt.Println("$ file print --name=NAME")
		kadem.PrintFile(takeName(words))
		break
	case "pin":
		//fmt.Println("$ file pin true/false")
		break
	case "list":
		kademlia.ListFiles()
		break
	default:
		fmt.Println("command not found")
	}

}
func processCommandInfo(words []string)  {
	if  len(words) != 1 || words[0] != "info" {
		fmt.Println("error INFO")
		return
	}
	fmt.Println("my info")
	//fmt.Println("***********")
}

func processCommandRoutingTable(words []string, kadem *kademlia.Kademlia){
	if  len(words) != 1 || words[0] != "routingTable" {
		fmt.Println("error routingTable")
		return
	}

	fmt.Println("my routing table")
	kadem.GetNetwork().GetMyRoutingTable().Print()
	//fmt.Println("***********")
}

func printHelp(){
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
	fmt.Println("")
	fmt.Println("Return the routing table")
	fmt.Println("$ lookup")
	fmt.Println("")
	fmt.Println("Createnew file in your node, random name")
	fmt.Println("$ file new")
	fmt.Println("")
	fmt.Println("Save file in the network")
	fmt.Println("$ file save --name=NAME")
	fmt.Println("")
	fmt.Println("Ping file")
	fmt.Println("$ file ping --name=NAME")
	fmt.Println("")
	fmt.Println("Find and Download file")
	fmt.Println("$ file take --name=NAME")
	fmt.Println("")
	fmt.Println("Remove file")
	fmt.Println("$ file rm --name=NAME")
	fmt.Println("")
	fmt.Println("Print data of a file")
	fmt.Println("$ file print --name=NAME")
	fmt.Println("")
	fmt.Println("Make file persistent or not")
	fmt.Println("$ file pin true/false")
	fmt.Println("")
	fmt.Println("List files in your node")
	fmt.Println("$ file list")
	fmt.Println("*****************************************")
}