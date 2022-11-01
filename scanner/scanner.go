// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package scanner

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"text/template"
	"unicode/utf8"
)

// Trace activates tracing for anything using the package. This is
// sometimes more convenient when an application uses the package but
// does not give access to the equivalent S.Trace property.
var Trace int

// ViewLen sets the number of bytes to view before eliding the rest.
var ViewLen = 20

var DefaultErrorMessage = `failed to scan`

// S (to avoid stuttering) implements a buffered data, non-linear,
// rune-centric, scanner with regular expression support.  Keep in mind
// that if and when you change the position (P) directly that rune (R)
// will not itself be updated as it is only updated by calling Scan.
// Often an update to the rune (R) as well would be inconsequential,
// even wasteful.  When less performant scanner operations are okay and
// or a higher level of abstraction allowed consider using the
// pegn.Scanner interface methods instead.
type S struct {
	B        []byte             // full buffer for lookahead or behind
	P        int                // index in B slice, points *after* R
	PP       int                // index of previous Scan, points *to* R
	R        rune               // last decoded, Scan updates, >1byte
	Trace    int                // activate trace log (>0)
	Errors   []error            // stack of errors in order
	Template *template.Template // for Report()
	NewLine  []string           // []string{"\r\n","\n"} by default
}

// New is a high-level scanner constructor and initializer that takes
// a single optional argument containing any valid Buffer() argument.
// Invalid arguments will fail (not fatal) with log output.
func New(args ...any) *S {
	s := new(S)
	if len(args) > 0 {
		s.Buffer(args[0])
	}
	return s
}

func (s *S) Bytes() []byte       { return s.B }
func (s *S) SetBytes(buf []byte) { s.B = buf }
func (s *S) Rune() rune          { return s.R }
func (s *S) SetRune(r rune)      { s.R = r }
func (s *S) Cur() int            { return s.P }
func (s *S) SetCur(p int)        { s.P = p }
func (s *S) Prev() int           { return s.PP }
func (s *S) SetPrev(p int)       { s.PP = p }

// Buffer sets the internal bytes buffer and initializes all internal
// pointers and state. This is useful when testing in order to buffer
// strings as well as content from any io.Reader.
func (s *S) Buffer(b any) {
	switch v := b.(type) {
	case string:
		s.B = []byte(v)
	case []byte:
		s.B = v
	case io.Reader:
		b, err := io.ReadAll(v)
		if err != nil {
			log.Printf("unable to read: %v", err)
			return
		}
		s.B = b
	}
	s.R = '\x00'
	s.P = 0
	s.PP = 0
}

const DefaultTemplate = `
{{- if .Errors -}}
	{{- range .Errors -}}
		error: {{.}}
	{{- end -}}
{{- else -}}
	{{- .Pos -}}
{{- end -}}
`

// Template is used by Report to log human-friendly scanner state
// information if the scan.R has not set its own Template.
var Template *template.Template

func init() {
	var err error
	Template, err = template.New("DefaultTemplate").Parse(DefaultTemplate)
	if err != nil {
		panic(err)
	}
}

// Position contains the human-friendly information about the position
// within a give text file. Note that all values begin with 1 and not
// 0.
type Position struct {
	Rune    rune // rune at this location
	BufByte int  // byte offset in file
	BufRune int  // rune offset in file
	Line    int  // line offset
	LByte   int  // line column byte offset
	LRune   int  // line column rune offset
}

// String fulfills the fmt.Stringer interface by printing
// the Position in a human-friendly way:
//
//   U+1F47F '👿' 1,3-5 (3-5)
//                | | |  | |
//             line | |  | overall byte offset
//   line rune offset |  overall rune offset
//     line byte offset
//
func (p Position) String() string {
	s := fmt.Sprintf(`%U %q %v,%v-%v (%v-%v)`,
		p.Rune, p.Rune,
		p.Line, p.LRune, p.LByte,
		p.BufRune, p.BufByte,
	)
	return s
}

// Print prints the cursor itself in String form. See String.
func (p Position) Print() { fmt.Println(p.String()) }

// Log calls log.Println on the cursor itself in String form. See String.
func (p Position) Log() { log.Println(p.String()) }

// Pos returns a human-friendly Position for the current location.
// When multiple positions are needed use Positions instead.

func (s S) Pos() Position { return s.Positions(s.P)[0] }

// Positions returns human-friendly Position information (which can easily
// be used to populate a text/template) for each raw byte offset (s.P).
// Only one pass through the buffer (s.B) is required to count lines and
// runes since the raw byte position (s.P) is frequently changed
// directly.  Therefore, when multiple positions are wanted, consider
// caching the raw byte positions (s.P) and calling Positions() once for
// all of them.
func (s S) Positions(p ...int) []Position {
	pos := make([]Position, len(p))

	if len(p) == 0 {
		return pos
	}

	if s.NewLine == nil {
		s.NewLine = []string{"\r\n", "\n"}
	}

	_rune, line, lbyte, lrune := 1, 1, 1, 1
	_s := S{B: s.B}
	//_s.Trace++

	for _s.Scan() {

		for _, nl := range s.NewLine {
			if _s.Is(nl) {
				line++
				_s.P += len(nl) - 1
				_rune += len(nl) - 1
				lbyte = 0
				lrune = 0
				continue
			}
		}

		for i, v := range p {
			if _s.P == v {
				pos[i] = Position{
					Rune:    _s.R,
					BufByte: _s.P,
					BufRune: _rune,
					Line:    line,
					LByte:   lbyte,
					LRune:   lrune,
				}
			}
		}

		rlen := len([]byte(string(s.R)))
		lbyte += rlen
		lrune++
		_rune++

	}

	return pos
}

// String implements fmt.Stringer with simply the position (P) and
// quoted rune (R) along with its Unicode. For printing more human
// friendly information about the current scanner state use Report.
func (s S) String() string {
	end := s.P + ViewLen
	elided := "..."
	if end > len(s.B) {
		end = len(s.B)
		elided = ""
	}
	return fmt.Sprintf("%v %q %q%v",
		s.P, s.R, s.B[s.P:end], elided)
}

// Print is shorthand for fmt.Println(s).
func (s S) Print() { fmt.Println(s) }

// Log is shorthand for log.Print(s).
func (s S) Log() { log.Println(s) }

// Scan decodes the next rune, setting it to R, and advances position
// (P) by the size of the rune (R) in bytes returning false then there
// is nothing left to scan. Only runes bigger than utf8.RuneSelf are
// decoded since most runes (ASCII) will usually be under this number.
func (s *S) Scan() bool {

	if s.P >= len(s.B) {
		return false
	}

	ln := 1
	r := rune(s.B[s.P])
	if r > utf8.RuneSelf {
		r, ln = utf8.DecodeRune(s.B[s.P:])
		if ln == 0 {
			return false
		}
	}

	s.PP = s.P
	s.P += ln
	s.R = r

	if s.Trace > 0 || Trace > 0 {
		s.Log()
	}

	return true
}

// Peek returns true if the passed string matches from current position
// in the buffer (s.P) forward. Returns false if the string
// would go beyond the length of buffer (len(s.B)).
func (s *S) Peek(a string) bool {
	if len(a)+s.P > len(s.B) {
		return false
	}
	if string(s.B[s.P:s.P+len(a)]) == a {
		return true
	}
	return false
}

// End returns true if scanner has nothing more to scan.
func (s *S) End() bool { return s.P == len(s.B) }

// Mark returns the main state values in order to jump Back() when
// required during other scan operations. Mark fulfills the pegn.Scanner
// interface.
func (s *S) Mark() (rune, int, int) { return s.R, s.P, s.PP }

// Back restores the main state of the scanner and fulfills the
// pegn.Scanner interface.
func (s *S) Back(r rune, p int, lp int) { s.R, s.P, s.PP = r, p, lp }

// Is returns true if the passed string matches the last scanned rune
// and the runes ahead matching the length of the string.  Returns false
// if the string would go beyond the length of buffer (len(s.B)).
func (s *S) Is(a string) bool {
	if len(a)+s.PP > len(s.B) {
		return false
	}

	if string(s.B[s.PP:s.PP+len(a)]) == a {
		return true
	}
	return false
}

// PeekMatch checks for a regular expression match at the current
// position in the buffer providing a mechanism for positive and
// negative lookahead expressions. It returns the length of the match.
// Successful matches might be zero (see regexp.Regexp.FindIndex).
// A negative value is returned if no match is found. Note that Go
// regular expressions now include the Unicode character classes (ex:
// \p{L}) that should be used over dated alternatives (ex: \w).
func (s *S) PeekMatch(re *regexp.Regexp) int {
	loc := re.FindIndex(s.B[s.P:])
	if loc == nil {
		return -1
	}
	if loc[0] == 0 {
		return loc[1]
	}
	return -1
}

// Match checks for a regular expression match at the last position in
// the buffer (s.PP) providing a mechanism for positive and negative
// lookahead expressions. It returns the length of the match.
// Successful matches might be zero (see regexp.Regexp.FindIndex).
// A negative value is returned if no match is found.  Note that Go
// regular expressions now include the Unicode character classes (ex:
// \p{L}) that should be used over dated alternatives (ex: \w).
func (s *S) Match(re *regexp.Regexp) int {
	loc := re.FindIndex(s.B[s.PP:])
	if loc == nil {
		return -1
	}
	if loc[0] == 0 {
		return loc[1]
	}
	return -1
}

// Report will fill in the s.Template (or scan.Template if not set) and
// log it to standard error. See the log package for removing prefixes
// and such. The DefaultTemplate is compiled at init() and assigned to
// the scan.Template global package variable. To silence reports
// developers may use the log package or simply ensure that both
// s.Template and scan.Template are nil.
func (s S) Report() {
	// TODO expand the s.Errors if no s.Position on first
	tmpl := s.Template
	if s.Template == nil {
		tmpl = Template
	}
	if tmpl == nil {
		return
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, s); err != nil {
		log.Println(err)
		return
	}
	log.Print(buf.String())
}

type Error struct {
	P   int      // can be left blank if Pos is defined
	Pos Position // can be left blank, Report will populate
	Msg string
}

func (e Error) Error() string {
	return fmt.Sprintf("%v at %v", e.Msg, e.Pos)
}

// Error adds an error to the Errors slice. Takes fmt.Sprintf() type
// arguments. The current position (s.Pos) is saved with the error.
// Since s.Pos scans to find the right location if there are multiple
// errors anticipated consider directly appending to Errors instead and
// only using the byte offset position (s.P). Report will detect if the
// first of Errors does not have a Position and will populate them
// efficiently before executing the template. For single errors, calling
// this method should be fine.
func (s *S) Error(a ...any) {
	msg := DefaultErrorMessage
	switch {
	case len(a) > 0:
		msg, _ = a[0].(string)
	case len(a) > 1:
		form, _ := a[0].(string)
		msg = fmt.Sprintf(form, a[1:]...)
	}
	s.Errors = append(s.Errors, Error{Pos: s.Pos(), Msg: msg})
}
