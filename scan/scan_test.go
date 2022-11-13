package scan_test

import (
	"fmt"

	"github.com/rwxrob/pegn/scan"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleField() {

	s := scanner.New(`fields don't have so-called spaces`)

	fmt.Println(scan.Field(s))
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(scan.Field(s))
	s.Print()

	// Output:
	// true
	// 's' 5-6 " don't hav"
	// ' ' 6-7 "don't have"
	// true
	// 't' 11-12 " have so-c"

}

func ExampleC_ws_scan() {

	s := scanner.New(`1 `)

	fmt.Println(scan.C_ws(s))
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(scan.C_ws(s))
	s.Print()

	// Output:
	// false
	// '\x00' 0-0 "1 "
	// '1' 0-1 " "
	// true
	// ' ' 1-2 ""

}
