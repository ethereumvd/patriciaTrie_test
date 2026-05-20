package main

import (
	"fmt"
)

//TODO - implement the remaining functions and embed the LRU pointers

// here a PatriciaTree struct has been used instead of a "Cache" struct
// appropriate changes has to be made in the actual codebase
// mutexes have also not been used for the sake of simplicity in this sample

// the cache (PatriciaTree) should support the following operations :-
//		checkInvariants() checks if some invariants have been violated
// 		evictOne() evicts a entry in the LRU cache if some invariant is violated
//  	Insert(key string, value any)
//  	eraseInternal(key string)  should erase a entry from the cache , lock not used
//  	eraseKeys(keys []string)
//  	Erase(key string) func to be used by other packages , in this fun use mutex to lock
//  	LookUp(key string) search for a key and update the LRU as it is most recently used , return val if found
//  	LookUpWithoutChangingOrder(key string) same but doesn't affect LRU
//  	UpdateWithoutChangingOrder(key string, val any) updates key is present otherwise return error
//

//helper functions for bit operations/manipulation

// returns the bit at a specified index of the given string
func bit(key string, index int) int {
	byteIndex := index / 8

	//pad with zeroes if requested index is out of bounds
	if byteIndex >= len(key) {
		return 0
	}

	bitIndex := 7 - (index % 8)
	return int((key[byteIndex] >> bitIndex) & 1)
}

// diffBit finds the first differing bit in the given keys
func diffBit(key1, key2 string) int {
	db := -1

	maxBits := 8 * (max(len(key1), len(key2)))

	for i := 0; i < maxBits; i++ {
		if bit(key1, i) != bit(key2, i) {
			db = i
			break
		}
	}

	return db
}

// TODO - embed LRU pointers in the Node and PatriciaTree struct
// by adding prev, next pointers in Node and head, tail pointers in PatriciaTree
type Node struct {
	isLeaf   bool
	bitIndex int
	key      string
	val      any
	left     *Node
	right    *Node
}

type PatriciaTree struct {
	root *Node
	//could add a mutex here maybe ??
}

func newPatriciaTree() *PatriciaTree {
	return &PatriciaTree{}
}

func (tr *PatriciaTree) insert(key string, val any) {

	//new tree case , create a new node and assign the key val
	if tr.root == nil {
		tr.root = &Node{
			isLeaf:   true,
			bitIndex: -1,
			key:      key,
			val:      val,
			left:     nil,
			right:    nil,
		}
		return
	}

	var par *Node
	curr := tr.root
	for !curr.isLeaf {
		par = curr
		if key[curr.bitIndex] == 1 {
			curr = curr.right
		} else {
			curr = curr.left
		}
	}

	//if we found the key is the same then just update the value
	if key == curr.key {
		curr.val = val
		return
	}

	//if the keys differ we have to create a new node and store the differing bit
	db := diffBit(key, curr.key)

	newLeaf := &Node{ //new entry in the tree
		isLeaf:   true,
		bitIndex: -1,
		key:      key,
		val:      val,
		left:     nil,
		right:    nil,
	}

	lft := &Node{}
	rt := &Node{}

	if curr.key[db] == 1 {
		rt = curr
		lft = newLeaf
	} else {
		lft = curr
		rt = newLeaf
	}

	newInternalNode := &Node{
		isLeaf:   false,
		bitIndex: -1,
		key:      "",
		val:      nil,
		left:     lft,
		right:    rt,
	}

	//edge case : if only one entry present i.e if root itself is a leaf
	if par == nil {
		tr.root = newInternalNode
		return
	}

	//check which side the leaf node was and put the new node accordingly
	if curr.key[par.bitIndex] == 1 {
		par.right = newInternalNode
	} else {
		par.left = newInternalNode
	}
}

func (tr *PatriciaTree) lookUp(key string) bool {

	//edge case if no entries present
	if tr.root == nil {
		return false
	}

	curr := tr.root
	for !curr.isLeaf {
		if key[curr.bitIndex] == 1 {
			curr = curr.right
		} else {
			curr = curr.left
		}
	}

	if curr.key == key {
		return true
	}
	return false
}

func main() {
	fmt.Println("hello")
}
