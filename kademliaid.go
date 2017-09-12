package d7024e

import (
	"encoding/hex"
	"math/rand"
)

const IDLength = 20

//KAdemliaId is an array  20 bytes, 160 bites
type KademliaID [IDLength]byte

//This functions returns a KademliaId
// params: data string
// hex.DecodeString(data) DecodeString returns the bytes represented by the hexadecimal string s.
// "48656c6c6f20476f7068657221" => Hello Gopher!
func NewKademliaID(data string) *KademliaID {
	decoded, _ := hex.DecodeString(data)

	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = decoded[i] //it returns the codes of the chars %!s(uint8=72) H
	}

	return &newKademliaID
}

//return an [160bits]array of random uint8=72
func NewRandomKademliaID() *KademliaID {
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

//Less is a method of KAdemliaId that returns a boolean if this currentNodeId is less than the KAdemliaId passed as param
//Params: KademliaID: other ID
func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	//This loop iterate over the two IDs if in some point they are different returns the comparation of CurrentNodeID[i] < PassedNodeID[i]
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	//If both nodes are the same it returns false
	return false
}

//Equal is a method of KAdemliaId that returns a boolean if this currentNodeId is equal than the KAdemliaId passed as param
func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

//REVIEW
//CalcDistance is a method of KAdemliaId
//params target is a KAdemliaID
//returns KAdemliaID
// ^ is a go operator (https://www.tutorialspoint.com/go/go_operators.htm)
//
// |---|---|-----|
// | p | q | p^q |
// |---|---|-----|
// | 0 | 0 |  0  |
// | 0 | 1 |  1  |
// | 1 | 1 |  0  |
// | 1 | 0 |  1  |
// |---|---|-----|

func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

// Srting is a method of KAdemliaId that returns the string of an Encoded KademliaID
//hex.EncodeToString(src) Hello -> 48656c6c6f
func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}
