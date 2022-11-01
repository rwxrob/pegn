/*

Package pegn provides utility functions for scanning and parsing
grammars defined in the Parsing Expression Grammar Notation (PEGN). (See
pegn.dev for details on the specification itself.) The package also
includes an abtract syntax tree structure that is consistent with the
model defined by PEGN with marshalling as both PEGN compressed and long
form JSON array formats as well as a more conventional JSON object node
tree.

Scanner Functions

Most of this package consists of scanner functions that might well have
been written using a pegn.Scan*(s) convention, but the sake of brevity
when writing recursive descent parsers the "Scan" has been dropped (ex:
pegn.ScanULETTER(s) -> pegn.ULETTER(s))

First Class Function Collection

Scanner Interface

*/
package pegn
