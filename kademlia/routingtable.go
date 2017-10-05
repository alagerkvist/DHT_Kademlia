package kademlia

import (
	"time"
)

const bucketSize = 10
const nb_task_managed = 100

type RoutingTable struct {
	me      Contact
	buckets [IDLength * 8]*bucket
	channelTasks chan Task
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
	bucket.AddContact(contact)
}

func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int, isForLookup bool) []Contact {
	var candidates ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target)

	bucket := routingTable.buckets[bucketIndex]
	if isForLookup{
		bucket.lastTimeVisited = time.Now().Local()
	}

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			if isForLookup{
				bucket.lastTimeVisited = time.Now().Local()
			}
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			if isForLookup{
				bucket.lastTimeVisited = time.Now().Local()
			}
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

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
	bucket.RemoveContact(contact)
}

func (routingTable *RoutingTable) GetMyContact() *Contact{
	return &routingTable.me
}


func (routingTable *RoutingTable) StartRoutingTableListener() {
	routingTable.channelTasks = make(chan Task, nb_task_managed)
	go routingTable.runWorker(routingTable.channelTasks)
}

func (routingTable *RoutingTable) createTask(idType int, responseChannel chan []Contact, contactRequested *Contact) *Task{
	var task Task = Task{idType, responseChannel, contactRequested}
	routingTable.channelTasks <- task
	return &task
}

