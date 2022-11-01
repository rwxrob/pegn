package scan_test

import (
	"fmt"

	"github.com/rwxrob/pegn/scan"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleSomeWS() {
	s := scanner.New(`1 `)
	//s.Trace++
	fmt.Println(scan.SomeWS(s))
	s.Print() // nothing advanced at all
	s.Scan()
	s.Print()
	fmt.Println(scan.SomeWS(s))
	s.Print()
	// Output:
	// false
	// 0 '\x00' "1 "
	// 1 '1' " "
	// true
	// 2 ' ' ""
}

func ExampleEndLine_carriage_Feed() {
	s := scanner.New("\r\n\r\n")
	s.Print()
	fmt.Println(scan.EndLine(s))
	s.Print()
	// Output:
	// 0 '\x00' "\r\n\r\n"
	// true
	// 2 '\n' "\r\n"
}

func ExampleEndLine_feed() {
	s := scanner.New("\n")
	s.Print()
	fmt.Println(scan.EndLine(s))
	s.Print()
	// Output:
	// 0 '\x00' "\n"
	// true
	// 1 '\n' ""
}

/*

func ExampleScanEndLine_carriage_Not_Enough() {
	s := new(scan.R)
	s.B = []byte("\r")
	s.Print()
	fmt.Println(basemd.ScanEndLine(s))
	s.Print()
	// Output:
	// 0 '\x00' "\r"
	// false
	// 0 '\x00' "\r"
}

/*
func ExampleScanEndParaline_Returns() {
	s := new(scan.R)
	s.B = []byte("\n\n")
	fmt.Println(basemd.ScanEOB(s))
	s.Print()
	// Output:
	// true
	// 2 '\n' ""
}

func ExampleScanEOB_carriage_and_Line_Returns() {
	s := new(scan.R)
	s.B = []byte("\r\n\r\n")
	fmt.Println(basemd.ScanEOB(s))
	s.Print()
	// Output:
	// true
	// 4 '\n' ""
}

func ExampleScanEOB_odd_Returns() {
	s := new(scan.R)
	s.B = []byte("\r\n\n")
	fmt.Println(basemd.ScanEOB(s))
	s.Print()
	// Output:
	// true
	// 3 '\n' ""
}

func ExampleScanEOB_extra_WS() {
	s := new(scan.R)
	s.B = []byte("   \r\n\r\n\r\n")
	fmt.Println(basemd.ScanEOB(s))
	s.Print()
	// Output:
	// true
	// 9 '\n' ""
}
*/
