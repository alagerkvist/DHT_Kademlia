package kademlia

import (
	"fmt"
	"sort"
)

/** A contact define a node:
* ID: the 160 bits address of the node
* Address: IP:PORT of the node
* distance: the distance from an other node
*/
type Contact struct {
	ID       *KademliaID
	Address  string
	distance *KademliaID
}

//this function returns a Contact with distance nil with the params ID and Address
func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil}
}

/** CalcDistance
* PARAM: contact: the contact to modify the distance parameter
*		 target: the kademlia ID to which we have to compute the distance
* Modify the field "distance" of contact to have the distance between contact and target
*/
func (contact *Contact) CalcDistance(target *KademliaID) {
	//this function uses the function CalcDistance of the KademliaId and
	//puts the returned value on the distance field of the contact
	contact.distance = contact.ID.CalcDistance(target)
}

/** Less
* PARAM: contact: the bucket to retrieve the contacts
*		 otherContact: the kademlia ID to which we have to compute the distance
*
* OUTPUT: If contact has a smaller distance than otherContact
*/func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.distance.Less(otherContact.distance)
}

// return the string of a contact
func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

//Array of the candidates of a Contact
type ContactCandidates struct {
	contacts []Contact
}

// Append is a method of ContactCandidates that appends the contacts passed into the contacts candidates
func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}

// GetContacts is a method os ContactCandidates that returns an array of contact with the length of the param count
func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	return candidates.contacts[:count]
}

// Sort is a method of ContactCandidates that sort the candidates with distance
func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}

// Len is a method of ContactCandidates that returns the length of the contactCandidates contacts
func (candidates *ContactCandidates) Len() int {
	return len(candidates.contacts)
}

// Swap is a method of ContactCandidates that swaps two contacts of the array, the params are the positions
func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.contacts[i], candidates.contacts[j] = candidates.contacts[j], candidates.contacts[i]
}

// Less is a method of ContactCandidates that return a boolean using the Contact.less method that uses the less mehod of KAdemliaId.
func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.contacts[i].Less(&candidates.contacts[j])
}
