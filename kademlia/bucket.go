package kademlia

import (
	"container/list"
	"fmt"
	"time"
)

/** Structure defining a bucket
* list: list of the contacts in the bucket
* cacheList: replacement list sorted by most recent view
* lastTimeVisited: last time that a lookup has been performed on this bucket
 */
type bucket struct {
	list *list.List
	cacheList *list.List
	lastTimeVisited time.Time
}


/** Initialize all the fields of the new bucket created
 */
func newBucket() *bucket {
	//constructor list of lists
	bucket := &bucket{}
	bucket.list = list.New()
	bucket.cacheList = list.New()
	return bucket
}

/** AddContact
* PARAM: bucket: the bucket to add a contact
*		 contact: the contact to add
* OUTPUT: if it is a new contact or nor
* Insert in the main list or in the waiting list the new contact depending if the main list is full or not
 */
 func (bucket *bucket) AddContact(contact Contact) bool{
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
				bucket.cacheList.PushFront(contact)
				return true
			} else{
				bucket.cacheList.MoveToFront(element)
			}

		}
	} else {
		// is the element exists we move the element to the front
		bucket.list.MoveToFront(element)
	}
	return false
}

/** RemoveContact
* PARAM: bucket: the bucket to remove a contact
*		 contact: the contact to remove
* Remove a contact from the main list bucket and add into it the first one of the waiting list
 */
func (bucket *bucket) RemoveContact(contact Contact){
	element := bucket.findElementInList(contact, bucket.list)
	if element == nil{
		fmt.Println("Error in deleting contact")
	} else {
		//Transfer the contact to the real bucket
		bucket.list.Remove(element)
		if bucket.cacheList.Len() > 0{
			bucket.list.PushFront(bucket.cacheList.Front())
			bucket.cacheList.Remove(bucket.cacheList.Front())
		}
	}
}


/** GetContactAndCalcDistance
* PARAM: bucket: the bucket to retrieve the contacts
*		 target: the kademlia ID to which we have to compute the distance

* OUTPUT: the list of contacts from the bucket with the distance parameter computed with the target contact
 */func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
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


/** GetContactAndCalcDistance
* PARAM: bucket: the bucket to check
*		 contact: the contact to find in the list
*		 ls: the list to retrieve the contact
* OUTPUT: The element found in the list
*/
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

//Print all the contacts of one bucket
func (bucket *bucket) Print() {
	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		fmt.Println(elt.Value.(Contact))
	}

}