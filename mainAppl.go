package main

import (
	//"time"
	/*"net"
	"fmt"
	"bufio"
	"github.com/golang/protobuf/proto"
	"log"
	"time"*/
	//"time"
	"./kademlia"
	"fmt"
	"bufio"
	"os"
	"strings"
	"net"

	"strconv"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	numberSrcNodes, _ :=  strconv.Atoi(os.Args[1])
	ip := getMyIp()
	port := "8080"

	var network kademlia.Network = kademlia.CreateRandomNetworks(numberSrcNodes, ip, port)
	kademlia.AddSourceNodes(&network, numberSrcNodes, ip, port)
	var kademlia *kademlia.Kademlia = &kademlia.Kademlia{&network}

	printHelp()

	for scanner.Scan() {
		processText(scanner.Text())
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


func processText(text string, kademlia *kademlia.Kademlia){
	//fmt.Println(text)
	//fmt.Println("hello")
	var words []string = strings.Split(text," ");

	var command = words[0]
	switch command {
	case "ping":
		processCommandPing(words)
		break
	case "info":
		processCommandInfo(words)
		break
	case "routingTable":
		processCommandRoutingTable(words)
		break
	case "help":
		break
	default:
		fmt.Println("command not found")
	}
	printHelp()
}

func processCommandPing(words []string){
	if  len(words) != 3 ||
		(words[1] != "--nodeID" && words[1] != "--nodeIP") {
		fmt.Println("error PING")
		return
	}
	fmt.Println("sending ping to ", words[2])
	fmt.Println("***********")
}

func processCommandInfo(words []string)  {
	if  len(words) != 1 || words[0] != "info" {
		fmt.Println("error INFO")
		return
	}
	fmt.Println("my info")
	fmt.Println("***********")
}

func processCommandRoutingTable(words []string)  {
	if  len(words) != 1 || words[0] != "routingTable" {
		fmt.Println("error routingTable")
		return
	}
	fmt.Println("my routing table")
	fmt.Println("***********")
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
	fmt.Println("*****************************************")
}
