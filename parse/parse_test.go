package parse_test

import (
	"fmt"

	"github.com/rwxrob/pegn/parse"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleC_ws_parse() {

	s := scanner.New(`1 `)

	fmt.Println(parse.C_ws(s))
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(parse.C_ws(s))
	s.Print()

	// Output:
	// <nil>
	// '\x00' 0-0 "1 "
	// '1' 0-1 " "
	// {"t":-77,"v":" "}
	// ' ' 1-2 ""

}
