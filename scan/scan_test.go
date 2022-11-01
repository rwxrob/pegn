package scan_test

import (
	"fmt"

	"github.com/rwxrob/pegn/scan"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleWS() {

	s := scanner.New(`1 `)

	fmt.Println(scan.WS(s))
	s.Print()

	// Output:
	// false
	// '\x00' 0-0 "1 "

}

func ExampleSomeWS() {

	s := scanner.New(`1 `)

	fmt.Println(scan.SomeWS(s))
	s.Print() // nothing advanced at all
	s.Scan()
	s.Print()
	fmt.Println(scan.SomeWS(s))
	s.Print()

	// Output:
	// false
	// '\x00' 0-0 "1 "
	// '1' 0-1 " "
	// false
	// ' ' 1-2 ""

}

func ExampleEndLine_carriage_Feed() {

	s := scanner.New("\r\n\r\n")

	s.Print()
	fmt.Println(scan.EndLine(s))
	s.Print()

	// Output:
	// '\x00' 0-0 "\r\n\r\n"
	// true
	// '\n' 1-2 "\r\n"

}

func ExampleEndLine_feed() {

	s := scanner.New("\n")
	s.Print()
	fmt.Println(scan.EndLine(s))
	s.Print()

	// Output:
	// '\x00' 0-0 "\n"
	// true
	// '\n' 0-1 ""

}

func ExampleScanEndLine_carriage_Not_Enough() {

	s := scanner.New("\r")

	s.Print()
	fmt.Println(scan.EndLine(s))
	s.Print()

	// Output:
	// '\x00' 0-0 "\r"
	// false
	// '\r' 0-1 ""

}

func ExampleScanEndPara_line_Returns() {

	s := scanner.New("\n\n")

	fmt.Println(scan.EndPara(s))
	s.Print()

	// Output:
	// true
	// '\n' 1-2 ""

}

func ExampleScanEndPara_carriage_and_Line_Returns() {

	s := scanner.New("\r\n\r\n")

	fmt.Println(scan.EndPara(s))
	s.Print()

	// Output:
	// true
	// '\n' 3-4 ""

}

func ExampleScanEndPara_odd_Returns() {

	s := scanner.New("\r\n\n")

	fmt.Println(scan.EndPara(s))
	s.Print()

	// Output:
	// true
	// '\n' 2-3 ""

}

func ExampleScanEndPara_extra_WS() {

	s := scanner.New("   \r\n\r\n\r\n")

	fmt.Println(scan.EndPara(s))
	s.Print()

	// Output:
	// true
	// '\n' 8-9 ""

}
