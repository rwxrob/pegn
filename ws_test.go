package pegn_test

import (
	"fmt"

	"github.com/rwxrob/pegn"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleCws() {

	s := scanner.New(`1 `)

	fmt.Println(pegn.Cws.Scan(s))
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(pegn.Cws.Scan(s))
	s.Print()

	// Output:
	// false
	// '\x00' 0-0 "1 "
	// '1' 0-1 " "
	// true
	// ' ' 1-2 ""

}
