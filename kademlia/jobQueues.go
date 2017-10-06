package kademlia

import (
"fmt"
)

type Task struct {
	idType int
	responseChan chan []Contact
	contactRequested *Contact
}

const lookUpContact = 0
const addContact = 1
const removeContact = 2
const getClosest = 3



func (routingTable *RoutingTable) runWorker(taskChannel <-chan Task){

	for {
		task := <-taskChannel

		//task.Print()

		switch task.idType {
		case lookUpContact:
			task.responseChan <- routingTable.FindClosestContacts(task.contactRequested.ID, bucketSize, true)
		case getClosest:
			task.responseChan <- routingTable.FindClosestContacts(task.contactRequested.ID, bucketSize, false)
		case addContact:
			fmt.Println("Contact to add")
			fmt.Println(*task.contactRequested)
			routingTable.AddContact(*task.contactRequested)
		case removeContact:
			routingTable.RemoveContact(*task.contactRequested)

		default:
			fmt.Printf("Error in task request")
		}
	}

}


func (task *Task) Print(){
	fmt.Print("* Task:")
	fmt.Println(task.idType)
	fmt.Println(task.contactRequested)
}
