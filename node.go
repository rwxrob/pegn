package pegn

import "encoding/json"

// Node is a typical node in a rooted node tree as required for any
// abstract syntax tree. Note that PEGN does not allow node attributes
// of any kind and that Nodes with Nodes under them must not also have
// a value (V).
//
// Minimal Design for Embedding
//
// The struct design is deliberately minimal with only Nodes, and String
// marshalling methods to ensure the least amount of conflict with
// potential embedded dependencies, which are encouraged to provide
// more involved handling of AST and other node trees when needed. Only
// unmarshaling methods are being considered at the moment (JSON, etc.).
//
// Unique Rule Type Association
//
// The Type integer often corresponds to the name (identifier) of the
// PEGN rule used to parse the Node, but not necessarily. For example,
// some rules are simply assertions that do not capture nodes. Others
// are not significant at all. Some don't even look at the content but
// examine the state of the scanner itself (Finished, etc.).
//
type Node struct {
	T int    // node type (linked to Rule.Type)
	V string // value (leaf nodes only)
	O *Node  // node over this node
	U *Node  // first node under this node (child)
	R *Node  // node to immediate right
}

// Nodes walks from the first node under (U) to the last returning a slice.
func (n Node) Nodes() []*Node {
	if n.U == nil {
		return nil
	}
	var nodes []*Node
	for cur := n.U; cur.R != nil; cur = cur.R {
		nodes = append(nodes, cur)
	}
	return nodes
}

// String fulfills the fmt.Stringer interface as JSON omitting any empty
// value (V) or slice of nodes under (N).
func (n Node) String() string {
	s := struct {
		T int
		V string  `json:",omitempty"`
		N []*Node `json:",omitempty"`
	}{n.T, n.V, n.Nodes()}
	byt, _ := json.Marshal(s)
	return string(byt)
}
