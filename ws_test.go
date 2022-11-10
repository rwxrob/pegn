package pegn_test

import (
	"fmt"

	"github.com/rwxrob/pegn"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleWhiteSpace_scan() {

	s := scanner.New(`1 `)

	fmt.Println(pegn.WhiteSpace.Scan(s))
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(pegn.WhiteSpace.Scan(s))
	s.Print()

	// Output:
	// false
	// '\x00' 0-0 "1 "
	// '1' 0-1 " "
	// true
	// ' ' 1-2 ""

}

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
