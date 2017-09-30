package kademlia

import (
	"container/list"
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
)

type Task struct {
	idType int
	doneRequest *bool
	contactRequested *Contact
	contactsReturn []Contact
}

const lookUpContact = 0
const addContact = 1
const removeContact = 2


type tasksList struct {
	list *list.List
}


func (routingTable *RoutingTable) runWorker(){

	for{
		if routingTable.listTasks.list.Len() != 0{
			task := routingTable.listTasks.list.Front().Value.(Task)

			switch task.idType {
			case lookUpContact:
				go PreFindClosestContacts(&task, routingTable)

			case addContact:
				go PreAddContact(&task, routingTable)

			case removeContact:
				go routingTable.RemoveContact(*task.contactRequested)

			default:
				fmt.Printf("Error in task request")
			}
		}
	}

}


func PreFindClosestContacts(task *Task, routingTable *RoutingTable){
	task.contactsReturn = routingTable.FindClosestContacts(task.contactRequested.ID, bucketSize)
	*task.doneRequest = true
}

func PreAddContact(task *Task, routingTable *RoutingTable){
	routingTable.AddContact(*task.contactRequested)
	//Can delete this request and pass to the next one
	routingTable.listTasks.list.Remove(routingTable.listTasks.list.Front())
}
