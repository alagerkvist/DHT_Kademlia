package kademlia

import (
	"container/list"
	"fmt"
	"time"
)

type bucket struct {
	list *list.List
	cacheList *list.List
	lastTimeVisited time.Time
}

func newBucket() *bucket {
	//constructor list of lists
	bucket := &bucket{}
	bucket.list = list.New()
	bucket.cacheList = list.New()
	return bucket
}

// AddContact is a method of bucket, params contact and it is void
func (bucket *bucket) AddContact(contact Contact) {
	//create a element variable
	 element := bucket.findElementInList(contact, bucket.list)

	//if the element doesn't exist
	if element == nil {
		if bucket.list.Len() < bucketSize {
			// if the list it is not full we push the new contact at the front
			bucket.list.PushFront(contact)
		} else{
			element = bucket.findElementInList(contact, bucket.cacheList)
			if element == nil{
				bucket.list.PushFront(contact)
			} else{
				bucket.list.MoveToFront(element)
			}

		}
	} else {
		// is the element exists we move the element to the front
		bucket.list.MoveToFront(element)
	}
}

func (bucket *bucket) RemoveContact(contact Contact){
	element := bucket.findElementInList(contact, bucket.list)
	if element == nil{
		fmt.Println("Error in deleting contact")
	} else {
		//Transfer the contact to the real bucket
		bucket.list.Remove(element)
		bucket.list.PushFront(bucket.cacheList.Front())
		bucket.cacheList.Remove(bucket.cacheList.Front())
	}
}

// GetContactAndCalcDistance is a method of bucket, params targer kademliaID and returns an Array of contacts
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	//create an array of contacts
	var contacts []Contact

	//this is to irerate over the list os the bucket.
	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		//now make operations over the contacts  REVIEW
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	//return the list of contacts
	return contacts
}

func(bucket *bucket) findElementInList(contact Contact, ls *list.List) *list.Element{
	//create a element variable
	var element *list.Element
	//this is to iterate over a list
	//list.front() returns the first element of a list

	for e := ls.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			//if the element is on the list, linear search , we save this element on the element variable
			element = e
		}
	}
	return element
}

//returns the len of the bucket.
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
