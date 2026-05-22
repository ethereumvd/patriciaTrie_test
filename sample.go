package main

import (
        "fmt"
        "strings"
)

//TODO - implement the remaining functions and embed the LRU pointers


// here a PatriciaTree struct has been used instead of a "Cache" struct
// appropriate changes has to be made in the actual codebase
// mutexes have also not been used for the sake of simplicity in this sample

// the cache (PatriciaTree) should support the following operations :-
//              checkInvariants() checks if some invariants have been violated
//              evictOne() evicts a entry in the LRU cache if some invariant is violated
//      Insert(key string, value any)
//      eraseInternal(key string)  should erase a entry from the cache , lock not used
//      eraseKeys(keys []string)
//      Erase(key string) func to be used by other packages , in this fun use mutex to lock
//      LookUp(key string) search for a key and update the LRU as it is most recently used , return val if found
//      LookUpWithoutChangingOrder(key string) same but doesn't affect LRU
//      UpdateWithoutChangingOrder(key string, val any) updates key is present otherwise return error
//              UpdateSize(key string, sizeDelta unit64) updates current size when entry's size has changed
//              EraseEntriesWithGivenPrefix(prefix string)

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
//Improvisation - can make a Node interface instead of a unified struct to save space something like :-
/*

type Node interface {
        isLeaf() bool
}

type InternalNode struct {
        bitIndex int
        left *Node
        right *Node
}

func (in *InternalNode) isLeaf() { return false }

type LeafNode struct{
        key string
        val any
}

func (ln *LeafNode) isLeaf() { return true }

*/

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
                        left:     ni,
                        right:    nil,
                }
                return
        }

        curr := tr.root
        for !curr.isLeaf {
                par = curr
                if bit(key, curr.bitIndex) == 1 {
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

        var par *Node
        curr = tr.root

        for !curr.isLeaf && curr.bitIndex < db {
                par = curr 

                if bit(key, curr.bitIndex) == 1 {
                        curr = curr.right
                } else {
                        curr = curr.left
                }
        }

        newLeaf := &Node{ //new entry in the tree
                isLeaf:   true,
                bitIndex: -1,
                key:      key,
                val:      val,
                left:     nil,
                right:    nil,
        }

        newInternalNode := &Node{
                isLeaf:   false,
                bitIndex: db,
                key:      "",
                val:      nil,
                left:     lft,
                right:    rt,
        }

        if bit(key, db) == 1 {
                newInternalNode.right = newLeaf
                newInternalNode.left = curr
        } else {
                newInternalNode.left = newLeaf
                newInternalNode.right = curr 
        }

        //edge case if no parent
        if par == nil {
                tr.root = newInternalNode
        } else {
                if par.left == curr {
                        par.left = newInternalNode
                } else {
                        par.right = newInternalNode
                }
        }

}

// as LRU has not been implemented yet, lookUp and lookUpWithoutChangingOrder are currently same
func (tr *PatriciaTree) lookUp(key string) (*Node, error) {

        err := fmt.Errorf("key ' %s ' not found", key)

        //edge case if no entries present
        if tr.root == nil {
                return nil, err
        }

        curr := tr.root
        for !curr.isLeaf {
                if bit(key, curr.bitIndex) == 1 {
                        curr = curr.right
                } else {
                        curr = curr.left
                }
        }

        if curr.key == key {
                return curr, nil
        }

        return nil, err
}

func (tr *PatriciaTree) UpdateWithoutChangingOrder(key string, newVal any) error {
        node, err := tr.lookUp(key)

        if err != nil {
                return err
        }

        node.val = newVal

        return nil
}

func (tr *PatriciaTree) eraseInternal(key string) error {
        _, err := tr.lookUp(key)
        if err != nil {
                return err
        }

        //edge case when root is itself the tree i.e no parent node
        //another edge case is when only one internal node is present i.e there is no grandparent of the leaf

        var par *Node
        var grPar *Node
        curr := tr.root

        for !curr.isLeaf {
                grPar = par
                par = curr
                if bit(key, curr.bitIndex) == 1 {
                        curr = curr.right
                } else {
                        curr = curr.left
                }
        }

        //if root itself is to be deleted
        if par == nil {
                tr.root = nil
                return nil
        }

        //if the leaf is at height 1 it may have a long sibling subtree
        //hence par is tr.root
        if grPar == nil {
                if par.left == curr {
                        par.key = par.left.key
                        par.val = par.left.val
                } else {
                        par.key = par.right.key
                        par.val = par.right.val
                }
                par.left = nil
                par.right = nil
                return nil
        }

        //normal case with grandparent
        //rearrange pointers
        if par.left == curr{
                if grPar.left == par {
                        grPar.left = par.left
                        par = nil
                } else {
                        grPar.right = par.left
                        par = nil
                }
        } else {
                if grPar.left == par {
                        grPar.left = par.right
                        par = nil
                } else {
                        grPar.right = par.right
                        par = nil
                }
        }

        return nil
}

func (tr *PatriciaTree) eraseKeys(keys []string) error {

        for _, key := range keys {
                err := tr.eraseInternal(key)
                if err != nil {
                        return err
                }
        }

        return nil
}

// this theoretically optimizes the current hashmap implementation's time complexity from O(number of keys) to O(max(length of keys))
// so in testcases where not a lot of entries are present but with long directory names (which is rare)
func (tr *PatriciaTree) eraseEntriesWithGivenPrefix(prefix string) error {

        //the first bitIndex >= len(prefix) * 8 will be where all nodes under it will have the desired prefix
        //we will need to delete such node
        //and attach it's sibling to it's grandparent eliminating the parent as well

        var par *Node
        var grPar *Node
        curr := tr.root

        for !curr.isLeaf && curr.bitIndex < len(prefix)*8 {
                grPar = par
                par = curr
                if bit(prefix, curr.bitIndex) == 1 {
                        curr = curr.right
                } else {
                        curr = curr.left
                }
        }

        //edge case if only root is there as a leaf we check if it has the prefix
        if par == nil {

                //if root itself is nil then no entries exist
                if tr.root == nil {
                        return fmt.Errorf("No entries exist")
                }

                if strings.HasPrefix(tr.root.key, prefix) {
                        tr.root = nil
                        return nil
                } else {
                        return fmt.Errorf("No entries with given prefix ' %s ' exist", prefix)
                }
        }

        //an edge case is if no entries of the prefix exist so we do a arbitrary traversal of the tree just to confirm

        search := curr
        for !search.isLeaf {
                search = search.right
        }
        if !strings.HasPrefix(search.key, prefix) {
                return fmt.Errorf("No entries with given prefix ' %s ' exist", prefix)
        }

        //edge case if height 1 binary tree i.e no grandparent
        if grPar == nil {
                if par.left == curr {
                        temp := tr.root
                        tr.root = tr.root.right
                        temp = nil
                } else {
                        temp := tr.root
                        tr.root = tr.root.left
                        temp = nil
                }
                return nil
        }

        //normal case have to rearrange pointers
        if(par == grPar.right) {
                if par.right == curr {
                        grpar.right = par.left
                } else {
                        grpar.right = par.right
                }
        } else {
                if par.right == curr {
                        grpar.left = par.left
                } else {
                        grpar.left = par.right
                }
        }
        par = nil

        return nil
}

func main() {
        fmt.Println("patriciatree test")
}

//BUG NOTICED :- In the case where we check for grPar != nil it is possible the other sibling of the node (leaf) to not be a leaf
