package kademlia

import (
	"container/list"
	"fmt"
)

const bucketSize = 10

type RoutingTable struct {
	me      Contact
	buckets [IDLength * 8]*bucket
	listTasks *tasksList
}

func NewRoutingTable(me Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = newBucket()
	}
	routingTable.me = me
	return routingTable
}

func (routingTable *RoutingTable) AddContact(contact Contact) {
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.mux.Lock()
	bucket.AddContact(contact)
	bucket.mux.Unlock()
}

func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) []Contact {
	var candidates ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target)

	bucket := routingTable.buckets[bucketIndex]
	bucket.mux.Lock()
	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket.mux.Unlock()
			bucket = routingTable.buckets[bucketIndex-i]
			bucket.mux.Lock()
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {
			bucket.mux.Unlock()
			bucket = routingTable.buckets[bucketIndex+i]
			bucket.mux.Lock()
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}
	bucket.mux.Unlock()

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetContacts(count)
}

func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(routingTable.me.ID)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return IDLength*8 - 1
}


func (routingTable *RoutingTable) RemoveContact(contact Contact){
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.mux.Lock()
	bucket.RemoveContact(contact)
	bucket.mux.Unlock()
}

func (routingTable *RoutingTable) GetMyContact() *Contact{
	return &routingTable.me
}


func (routingTable *RoutingTable) StartRoutingTableListener() {
	routingTable.listTasks = &tasksList{}
	routingTable.listTasks.list = list.New()
	routingTable.runWorker()
}

func (routingTable *RoutingTable) createTask(idType int, doneRequest *bool, contactRequested *Contact, contactsReturn []Contact) *Task{
	var task *Task = &Task{idType, doneRequest, contactRequested, contactsReturn}
	fmt.Println(task)
	routingTable.listTasks.list.PushBack(task)
	return task
}

func (routingTable *RoutingTable) lookUpContactRequest(kademliaId *KademliaID) []Contact{
	receivedContact := make([]Contact, bucketSize)
	endRequest := false
	for !endRequest{
		routingTable.createTask(lookUpContact, &endRequest, &Contact{kademliaId, "", nil}, receivedContact)
	}

	return receivedContact
}
