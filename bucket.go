package d7024e

import (
	"container/list"
)

type bucket struct {
	list *list.List
}

func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

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

func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
