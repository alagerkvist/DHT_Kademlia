package kademlia

import (
	"container/list"
)

type bucket struct {
	list *list.List
}

func newBucket() *bucket {
	//constructor list of lists
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

// AddContact is a method of bucket, params contact and it is void
func (bucket *bucket) AddContact(contact Contact) {
	//create a element variable
	var element *list.Element
	//this is to iterate over a list
	//list.front() returns the first element of a list

	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			//if the element is on the list, linear search , we save this element on the element variable
			element = e
		}
	}

	//if the element doesn't exist
	if element == nil {
		if bucket.list.Len() < bucketSize {
			// if the list it is not full we push the new contact at the front
			bucket.list.PushFront(contact)
		}
	} else {
		// is the element exists we move the element to the front
		bucket.list.MoveToFront(element)
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

//returns the len of the bucket.
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
