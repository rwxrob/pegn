package scan_test

import (
	"fmt"

	"github.com/rwxrob/pegn/scan"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleWhiteSpace_scan() {

	s := scanner.New(`1 `)

	fmt.Println(scan.WhiteSpace(s))
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(scan.WhiteSpace(s))
	s.Print()

	// Output:
	// false
	// '\x00' 0-0 "1 "
	// '1' 0-1 " "
	// true
	// ' ' 1-2 ""

}
