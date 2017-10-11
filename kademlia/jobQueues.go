package kademlia

import (
"fmt"
)

/** Task: structure that define a task to perform for this go routine
*/
type Task struct {
	idType int
	responseChan chan []Contact
	contactRequested *Contact
}

const lookUpContact = 0
const addContact = 1
const removeContact = 2
const getClosest = 3


/** runWorker
* PARAM: routingTable
*		 taskChannel: the channel to receive the tasks
* Perform the tasks one by one
*/
func (kademlia *Kademlia) runWorker(taskChannel <-chan Task){

	for {
		task := <-taskChannel

		switch task.idType {
		case lookUpContact:
			task.responseChan <- kademlia.network.myRoutingTable.FindClosestContacts(task.contactRequested.ID, bucketSize, true)
		case getClosest:
			task.responseChan <-  kademlia.network.myRoutingTable.FindClosestContacts(task.contactRequested.ID, bucketSize, false)
		case addContact:
			if kademlia.network.myRoutingTable.AddContact(*task.contactRequested){
				//me := kademlia.network.myRoutingTable.me
				/*for k := range kademlia.network.FileManager.filesStored {
					closestNode := kademlia.network.myRoutingTable.FindClosestContacts(NewKademliaID(k), 1, false)
					me.CalcDistance(NewKademliaID(k))
					if me.Less(&closestNode[0]){
						data := kademlia.network.FileManager.readData(k)
						base64Data := base64.StdEncoding.EncodeToString(data[:])
						kademlia.network.marshalStore(k, base64Data, task.contactRequested)
					}
				}*/
			}
		case removeContact:
			kademlia.network.myRoutingTable.RemoveContact(*task.contactRequested)

		default:
			fmt.Printf("Error in task request")
		}
	}

}

// Print a Taskrun the task
func (task *Task) Print(){
	fmt.Print("* Task:")
	fmt.Println(task.idType)
	fmt.Println(task.contactRequested)
}
