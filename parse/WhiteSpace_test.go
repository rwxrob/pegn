package parse_test

import (
	"fmt"

	"github.com/rwxrob/pegn"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleWhiteSpace_parse() {

	s := scanner.New(`1 `)

	fmt.Println(pegn.WhiteSpace.Parse(s))
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(pegn.WhiteSpace.Parse(s))
	s.Print()

	// Output:
	// <nil>
	// '\x00' 0-0 "1 "
	// '1' 0-1 " "
	// {"T":1,"V":" "}
	// ' ' 1-2 ""

}
