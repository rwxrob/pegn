// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/rwxrob/pegn/qstack"
)

// Node is an implementation of a "node" from traditional rooted node
// tree data structure. A Node may have other nodes under then (which
// some call a "parent" (P)) or can be an "edge" or "leaf" with a value
// instead. While there is nothing preventing a Node from having both
// a value and other nodes under it, such use is unsupported by the
// MarshalJSON/UnmarshalJSON methods. All nodes have a specific integer
// type (T).
type Node struct {
	T     int    `json:"T"`          // type
	V     string `json:",omitempty"` // value
	P     *Node  `json:"-"`          // up/parent
	Count int    `json:"-"`          // node count

	left  *Node
	right *Node
	first *Node
	last  *Node
}

// Init resets the node to its empty/zero state as if just created for
// the first time.
func (n *Node) Init() {
	n.T = 0
	n.V = ""
	n.first = nil
	n.last = nil
	n.left = nil
	n.right = nil
}

// Nodes returns all the nodes under this Node. Prefer checking
// Count when the values are not needed.
func (n *Node) Nodes() []*Node {
	if n.first == nil {
		return nil
	}
	cur := n.first
	list := []*Node{cur}
	for {
		cur = cur.right
		if cur == nil {
			break
		}
		list = append(list, cur)
	}
	return list
}

// --------------------------------------------------------------------

// Add creates a new Node with type and value under and returns. It also
// updates Count.
func (n *Node) Add(t int, v string) *Node {
	u := new(Node)
	u.T = t
	u.V = v
	u.P = n
	n.Append(u)
	return u
}

// Cut removes a Node from under the one above it and returns.
func (n *Node) Cut() *Node {
	if n.left != nil {
		n.left.right = n.right
	}
	if n.right != nil {
		n.right.left = n.left
	}
	if n.P != nil {
		n.P.Count--
		if n == n.P.first {
			n.P.first = n.right
		}
		if n == n.P.last {
			n.P.last = n.left
		}
	}
	n.P = nil
	n.left = nil
	n.right = nil
	return n
}

// Take moves all nodes from another under itself.
func (n *Node) Take(from *Node) {
	if from.first == nil {
		return
	}
	if n.first == nil {
		n.first = from.first
		n.last = from.last
		n.Count = from.Count
	} else {
		n.Count += from.Count
		n.last.right = from.first
		from.first.left = n.last
		n.last = from.last
	}
	from.Count = 0
	from.first = nil
	from.last = nil
}

// Append adds an existing Node under this one as if Add had been
// called.
func (n *Node) Append(u *Node) {
	n.Count++
	if n.first == nil {
		n.first = u
		n.last = u
		return
	}
	n.last.right = u
	u.left = n.last
	n.last = u
}

// Morph initializes the node with Init and then sets it's value (V) and
// type (T) and all of its attachment references to those of the Node
// passed thereby preserving the Node reference of this method's
// receiver.
func (n *Node) Morph(c *Node) {
	n.Init()
	n.T = c.T
	n.V = c.V
	n.P = c.P
	n.left = c.left
	n.right = c.right
	n.first = c.first
	n.last = c.last
	n.Count = c.Count
}

// Refs returns the internal pointers as a string for visualization
// mostly during debugging. See LogRefs.
func (n *Node) Refs() string {
	return fmt.Sprintf(`self:  %p parent: %p
left:  %-12p right:  %p
first: %-12p last:   %p`, n, n.P, n.left, n.right, n.first, n.last)
}

// Copy returns a duplicate of the Node and all its relations. Values
// are copied using simple assignment. Copy is useful for preserving
// state in order to revert a Node or to allow independent processing
// with concurrency on individual copies. Note that Node[<ref>] types
// will not produce deep copies of values.
func (n *Node) Copy() *Node {
	clones := map[*Node]*Node{}
	list := qstack.New[*Node]()
	list.Unshift(n)
	for list.Len > 0 {
		cur := list.Shift()
		list.Unshift(cur.Nodes()...)
		c := *cur
		clones[cur] = &c
	}
	for _, clone := range clones {
		clone.P = clones[clone.P]
		clone.left = clones[clone.left]
		clone.right = clones[clone.right]
		clone.first = clones[clone.first]
		clone.last = clones[clone.last]
	}
	return clones[n]
}

// ------------------------------- Walk --------------------------------

// WalkLevels will pass each Node in the tree to the given function
// traversing in a synchronous, breadth-first, leveler way. The function
// passed may be a closure containing variables, contexts, or a channel
// outside of its own scope to be updated for each visit. This method
// uses functional recursion which may have some limitations depending
// on the depth of node trees required.
func (n *Node) WalkLevels(do func(n *Node)) {
	list := qstack.New[*Node]()
	list.Unshift(n)
	for list.Len > 0 {
		cur := list.Shift()
		list.Push(cur.Nodes()...)
		do(cur)
	}
}

// WalkDeepPre will pass each Node in the tree to the given function
// traversing in a synchronous, depth-first, preorder way. The function
// passed may be a closure containing variables, contexts, or a channel
// outside of its own scope to be updated for each visit. This method
// uses functional recursion which may have some limitations depending
// on the depth of node trees required.
func (n *Node) WalkDeepPre(do func(n *Node)) {
	list := qstack.New[*Node]()
	list.Unshift(n)
	for list.Len > 0 {
		cur := list.Shift()
		list.Unshift(cur.Nodes()...)
		do(cur)
	}
}

// ------------------------------ Printer -----------------------------
// just for marshaling
type jsnode struct {
	T int     `json:"T"`
	V string  `json:"V,omitempty"`
	N []*Node `json:"N,omitempty"`
}

// MarshalJSON fulfills the json.Marshaler interface by first creating
// a copy of itself with the Nodes has N and then marshaling with an
// encoder that has had SetEscapeHTML set to false trims the extraneous
// newline added by json.Encoder.Encode. See String, Log, and Print as well.
func (s Node) MarshalJSON() ([]byte, error) {
	n := new(jsnode)
	n.T = s.T
	n.V = s.V
	n.N = s.Nodes()
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(n)
	byt := buf.Bytes()
	return byt[:len(byt)-1], err
}

// String returns the MarshalJSON version or the string "null" if an
// error occurred. An error is also logged with log.Print. No additional
// line return is added.
func (s Node) String() string {
	byt, err := s.MarshalJSON()
	if err != nil {
		log.Println(err)
		return `null`
	}
	return string(byt)
}

// Print uses fmt.Print to print.
func (s Node) Print() { fmt.Print(s.String()) }

// Println uses fmt.Println to print.
func (s Node) Println() { fmt.Println(s.String()) }

// Log uses log.Print to print.
func (s Node) Log() { log.Print(s.String()) }
