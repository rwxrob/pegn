package scan_test

import (
	"fmt"

	"github.com/rwxrob/pegn/scan"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleField_scan() {

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
