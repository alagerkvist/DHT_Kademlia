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
	"fmt"
	"bufio"
	"os"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("hello, I am a Kademlia node.")
	printHelp()

	for scanner.Scan() {
		processText(scanner.Text())
	}
}


func processText(text string){
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
