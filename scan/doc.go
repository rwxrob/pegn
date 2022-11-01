/*

Package scan contains a large library of scan functions for use by those
developing PEG-centric recursive descent parsers.

Simple and Adequately Performant

A balance between performance and simplicity has been the primary goal
with this package. The thought is that the time saved quickly writing
recursive descent parser prototypes using this package can be applied
later to convert these into a more performant forms having done the hard
initial work of figuring out the grammar notation and parsing algorithm
required.

These scan functions use of a pegn.Scanner with a cursor that contains
the last rune scanned and both it's beginning and end index in the bytes
buffer. Using a cursor in this way adds several steps of indirection
beyond the function calls themselves, but this makes for scanners that
are very easy to quickly write, trace, test, and maintain.

The simple choice to use a functional scanner approach over a single
loop is also less performant, but, again, much easier to write, maintain,
and share, which is also why the functions take a pegn.Scanner interface
as their only parameter (pegn.ScanFunc).

Functions that fulfill the pegn.ScanFunc interface are guaranteed to be
100% compatible with any pegn.Scanner implementation, forever. This
promotes community contribution of pegn.ScanFuncs for reuse as imported
and first-class functions, whatever the application may be.

*/
package scan
