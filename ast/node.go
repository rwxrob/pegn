package ast

import (
	"encoding/json"
)

type Node struct {
	T int    // node type (linked to Rule.Type)
	V string // value (leaf nodes only)
	O *Node  // node over this node
	U *Node  // first node under this node (child)
	R *Node  // node to immediate right
}

// Nodes walks from the first node under (U) to the last returning a slice.
func (n *Node) Nodes() []*Node {
	if n.U == nil {
		return nil
	}
	var nodes []*Node
	for cur := n.U; cur.R != nil; cur = cur.R {
		nodes = append(nodes, cur)
	}
	return []*Node(nodes)
}

// String fulfills the fmt.Stringer interface as JSON omitting any empty
// value (V) or slice of nodes under (N).
func (n *Node) String() string {
	s := struct {
		T int     `json:"t"`
		V string  `json:"v,omitempty"`
		N []*Node `json:"n,omitempty"`
	}{n.T, n.V, n.Nodes()}
	// FIXME do the required validation of V and N
	byt, _ := json.Marshal(s)
	// FIXME convert to "error: " string if error (per pegn.Node interface)
	return string(byt)
}
