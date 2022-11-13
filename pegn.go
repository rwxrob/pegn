/*

Package pegn provides tools for working with grammars written in Parsing
Expression Grammar Notation (Rob Muhlestein's combination of ABNF with
Bryan Ford's original PEG design as found in the examples in his paper).

Tools

	  * scan/     - scan functions:  func(s Scanner) bool
		* parse/    - parse functions: func(s Scanner) Node
		* is/       - class functions: func(a rune) bool
    * tk/       - token constants
    * types.go  - type constants
		* docs/     - multi-lingual rule documentation
    * scanner/  - pegn.Scanner implementation (DefaultScanner)
    * curs/     - curs.R struct for position within bytes buffer
    * cmd.go    - Bonzai stateful command tree composition
    * cmd/pegn  - helper utility with code generation

Types of Rules

A PEGN rule is what appears on the right of the arrow (<-). There are
two subtypes of a PEGN rule: class rules and a token rules. A class is
a set of runes. A token is a specific rune or string constant.

All rules must have a unique integer constant type.

PEGN reserved integer types (types.go) are guaranteed to never change
and always be negative. This frees grammar developers to use positive
integers. On day, the greater PEGN community may wish to organize range
reservations for different common grammar rules to maintain grammar and
AST interoperability.

No guarantee is made about what numeric range a rule, class, or token
integer will be only that any specific type integer will never be reused
for a different type. Use of such value ranges is strongly discouraged
over creation of proper range maps (as is used in the unicode package).

The integer 0 is reserved as Untyped.

Every PEGN rule must have a ScanFunc (scan) and a ParseFunc (parse).
Class rules must also have a ClassFunc (is). Token rules must also have
a token constant (tk).

Balanced Simplicity and Performance

A balance between performance and simplicity has been the primary goal
with this package. The thought is that the time saved quickly writing
recursive descent parser prototypes using this package can be applied
later to convert these into a more performant forms having done the hard
initial work of figuring out the grammar notation and parsing algorithm
required.

*/
package pegn

import (
	"go/ast"

	"github.com/rwxrob/pegn/gr"
	"github.com/rwxrob/pegn/scanner"
)

var DefaultScanner = scanner.New()

// Scan uses the DefaultScanner to load whatever input the Scanner
// implementation accepts to its Buffer method and then parses the PEGN
// spec string into a Grammar and delegates Scan to it.
func Scan(in any, spec string) (bool, []error) { return gr.PEGN.Scan(in, spec) }

// Parse does the same thing as Scan but produces a tree of parsed nodes
// instead of just validating.
func Parse(in any, spec string) (*ast.Node, []error) { return gr.PEGN.Parse(in, spec) }

func ScanClass(s Scanner, f ClassFunc) bool {
	// TODO
	return false
}

func ParseClass(s Scanner, f ClassFunc) *ast.Node {
	// TODO
	return nil
}
