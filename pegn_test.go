package pegn_test

import (
	"fmt"

	"github.com/rwxrob/pegn"
)

func ExampleNode_nodes() {
	n := new(pegn.Node)
	u1 := new(pegn.Node)
	u2 := new(pegn.Node)
	n.U = u1
	n.U.R = u2
	fmt.Println(n.Nodes())
	// Output:
	// [{"T":0}]
}
