package pegn_test

import (
	"fmt"

	"github.com/rwxrob/pegn"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleUprint_scan() {

	s := scanner.New(`d t`)

	// note space is included, tabs and returns not

	fmt.Println(pegn.Uprint.Scan(s))
	s.Print()
	fmt.Println(pegn.Uprint.Scan(s))
	s.Print()
	fmt.Println(pegn.Uprint.Scan(s))
	s.Print()

	// Output:
	// true
	// 'd' 0-1 " t"
	// true
	// ' ' 1-2 "t"
	// true
	// 't' 2-3 ""

}

func ExampleUprint_parse() {

	s := scanner.New(`d`)

	fmt.Println(pegn.Field.Parse(s))
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(pegn.Field.Parse(s))
	s.Print()

	// Output:
	// {"T":2,"V":"d"}
	// 'd' 0-1 ""
	// 'd' 0-1 ""
	// <nil>
	// 'd' 0-1 ""

}
