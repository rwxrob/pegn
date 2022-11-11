package pegn_test

import (
	"fmt"

	"github.com/rwxrob/pegn"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleField_scan() {

	s := scanner.New(`fields don't have so-called spaces`)

	fmt.Println(pegn.Field.Scan(s))
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(pegn.Field.Scan(s))
	s.Print()

	// Output:
	// true
	// 's' 5-6 " don't hav"
	// ' ' 6-7 "don't have"
	// true
	// 't' 11-12 " have so-c"

}

func ExampleField_parse() {

	s := scanner.New(`fields don't have so-called spaces`)

	fmt.Println(pegn.Field.Parse(s))
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(pegn.Field.Parse(s))
	s.Print()

	// Output:
	// {"T":2,"V":"fields"}
	// 's' 5-6 " don't hav"
	// ' ' 6-7 "don't have"
	// {"T":2,"V":"don't"}
	// 't' 11-12 " have so-c"

}
