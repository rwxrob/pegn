package ast_test

import (
	"fmt"

	"github.com/rwxrob/pegn/ast"
)

func ExampleNode_Init() {

	// create and print a brand new one
	n := new(ast.Node)
	n.Println()

	// add something to it
	n.V = `something`
	n.Println()

	// initialize it back to empty zero value
	n.Init()
	n.Println()
	// Output:
	// {"T":0}
	// {"T":0,"V":"something"}
	// {"T":0}
}

func ExampleNode_String() {
	t := new(ast.Node)
	t.V = `<foo>` // <> not be escaped by encoding/json
	t.Print()
	// Output:
	// {"T":0,"V":"<foo>"}
}

func ExampleNode_properties() {

	// Nodes have these properties updating every time
	// their state is changed so that queries need not
	// do the checks again later.

	// initial state
	n := new(ast.Node)
	fmt.Printf("n: %v %q %v\n", n.P == nil, n.V, n.Count)
	u := n.Add(1, "")
	fmt.Printf("n: %v %q %v\n", n.P == nil, n.V, n.Count)
	fmt.Printf("u: %v %q %v\n", u.P == nil, u.V, u.Count)

	// make an edge node
	u.V = "something"

	// break edge by forcing it to have nodes and a value (discouraged)
	u.Add(9001, "muhaha")
	fmt.Printf("u: %v %q %v\n", u.P == nil, u.V, u.Count)

	// Output:
	// n: true "" 0
	// n: true "" 1
	// u: false "" 0
	// u: false "something" 1

}

func ExampleNode_Nodes() {

	n := new(ast.Node)
	n.Add(1, "")
	n.Add(2, "")
	fmt.Println(n.Nodes(), n.Count)

	// and another added under it
	m := n.Add(3, "")
	m.Add(3, "")
	m.Add(3, "")
	fmt.Println(m.Nodes(), m.Count)
	fmt.Println(n.Nodes(), n.Count)

	// Output:
	// [{"T":1} {"T":2}] 2
	// [{"T":3} {"T":3}] 2
	// [{"T":1} {"T":2} {"T":3,"N":[{"T":3},{"T":3}]}] 3

}

func ExampleNode_Cut_middle() {
	n := new(ast.Node)
	n.Add(1, "")
	c := n.Add(2, "")
	n.Add(3, "")
	n.Println()
	fmt.Println(n.Count)
	x := c.Cut()
	n.Println()
	fmt.Println(n.Count)
	x.Println()
	// Output:
	// {"T":0,"N":[{"T":1},{"T":2},{"T":3}]}
	// 3
	// {"T":0,"N":[{"T":1},{"T":3}]}
	// 2
	// {"T":2}
}

func ExampleNode_Cut_first() {
	n := new(ast.Node)
	c := n.Add(1, "")
	n.Add(2, "")
	n.Add(3, "")
	n.Println()
	x := c.Cut()
	n.Println()
	x.Println()
	// Output:
	// {"T":0,"N":[{"T":1},{"T":2},{"T":3}]}
	// {"T":0,"N":[{"T":2},{"T":3}]}
	// {"T":1}
}

func ExampleNode_Cut_last() {
	n := new(ast.Node)
	n.Add(1, "")
	n.Add(2, "")
	c := n.Add(3, "")
	n.Println()
	x := c.Cut()
	n.Println()
	x.Println()
	// Output:
	// {"T":0,"N":[{"T":1},{"T":2},{"T":3}]}
	// {"T":0,"N":[{"T":1},{"T":2}]}
	// {"T":3}
}

func ExampleNode_Take() {

	// build up the first
	n := new(ast.Node)
	n.T = 10
	n.Add(1, "")
	n.Add(2, "")
	n.Add(3, "")
	n.Println()
	fmt.Println(n.Count)

	// now take them over

	m := new(ast.Node)
	m.T = 20
	m.Println()
	fmt.Println(m.Count)
	m.Take(n)
	m.Println()
	fmt.Println(m.Count)
	n.Println()
	fmt.Println(n.Count)

	// Output:
	// {"T":10,"N":[{"T":1},{"T":2},{"T":3}]}
	// 3
	// {"T":20}
	// 0
	// {"T":20,"N":[{"T":1},{"T":2},{"T":3}]}
	// 3
	// {"T":10}
	// 0

}

func ExampleNode_WalkLevels() {
	n := new(ast.Node)
	n.Add(1, "").Add(11, "")
	n.Add(2, "").Add(22, "")
	n.Add(3, "").Add(33, "")
	n.WalkLevels(func(c *ast.Node) { fmt.Print(c.T, " ") })
	// Output:
	// 0 1 2 3 11 22 33
}

func ExampleNode_WalkDeepPre() {
	n := new(ast.Node)
	n.Add(1, "").Add(11, "")
	n.Add(2, "").Add(22, "")
	n.Add(3, "").Add(33, "")
	n.WalkDeepPre(func(c *ast.Node) { fmt.Print(c.T, " ") })
	// Output:
	// 0 1 11 2 22 3 33
}

func ExampleNode_Morph() {
	n := new(ast.Node)
	n.Add(2, "some")
	m := new(ast.Node)
	m.Morph(n)
	n.Println()
	m.Println()
	// Output:
	// {"T":0,"N":[{"T":2,"V":"some"}]}
	// {"T":0,"N":[{"T":2,"V":"some"}]}
}

func ExampleNode_Copy() {
	n := new(ast.Node)
	n.Add(2, "some")

	c := n.Copy()
	c.Add(3, "new").Add(4, "deep")

	// 	log.Println("Original -------------------------")
	// 	n.LogRefs()
	// 	log.Println("Copy -----------------------------")
	// 	c.LogRefs()

	fmt.Println(&n != &c)
	n.Println()
	c.Println()

	// Output:
	// true
	// {"T":0,"N":[{"T":2,"V":"some"}]}
	// {"T":0,"N":[{"T":2,"V":"some"},{"T":3,"V":"new","N":[{"T":4,"V":"deep"}]}]}

}
