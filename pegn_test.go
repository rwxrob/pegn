package pegn_test

import (
	"fmt"

	"github.com/rwxrob/pegn"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleIs_ws() {

	fmt.Println(pegn.Is_ws(' '))
	fmt.Println(pegn.Is_ws('\r'))
	fmt.Println(pegn.Is_ws('\t'))
	fmt.Println(pegn.Is_ws('\n'))
	fmt.Println(pegn.Is_ws('\x00'))
	fmt.Println(pegn.Is_ws('1'))

	// Output:
	// true
	// true
	// true
	// true
	// false
	// false

}

func ExampleScan_ws() {

	s := scanner.New(`1 `)

	fmt.Println(pegn.Scan_ws(s, nil))
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(pegn.Scan_ws(s, nil))
	s.Print()

	// Output:
	// false
	// '\x00' 0-0 "1 "
	// '1' 0-1 " "
	// true
	// ' ' 1-2 ""

}

func ExampleParse_ws() {

	s := scanner.New(`1 `)

	fmt.Println(pegn.Parse_ws(s))
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(pegn.Parse_ws(s))
	s.Print()

	// Output:
	// <nil>
	// '\x00' 0-0 "1 "
	// '1' 0-1 " "
	// {"T":-1,"V":" "}
	// ' ' 1-2 ""

}
