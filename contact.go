package d7024e

import (
	"fmt"
	"sort"
)

//A Contact is one node it has one ID KademliaID, one address and one distance
type Contact struct {
	ID       *KademliaID
	Address  string
	distance *KademliaID
}

//this function returns a Contact with distance nil with the params ID and Address
func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil}
}

//CalcDistance is a method of Contact that has one KademliaId as param and it is void.
func (contact *Contact) CalcDistance(target *KademliaID) {
	//this function uses the function CalcDistance of the KademliaId and
	//puts the returned value on the distance field of the contact
	contact.distance = contact.ID.CalcDistance(target)
}

//Less is a method of Contact that return the bool value of the KAdemliaId.less() using as param otherContact.
func (contact *Contact) Less(otherContact *Contact) bool {
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

// Sort is a method of ContactCandidates that sort the candidates
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
