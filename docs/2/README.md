# Design considerations

* [Use `C_` prefix for classes to preserve identifiers](../3?L)
* [Simplest rune pegn.Scanner interface](../4?L)
* [Decouple meta-data association](../5?L)
* [All scanner functions advance the scanner themselves](../6?L)

Since it is not always clear how much of the bytes buffer will be needed, the task of advancing the scanner (with `s.Scan()` or however) is up to each scan function.

* **Scan functions must not change state on fail**

The scanner performing the scan must be in exactly the same state (rune,
byte pointer, previous byte pointer, etc.) if the scan fails (returning
`false`). The only exception is pushing one or more errors into the scanner's memory.

* **Scan functions must advance on success**

If a scanner function successfully matches it must leave the scanner
pointing to the next byte in the bytes buffer (not yet scanned) ready
for the next scan function.

* **Scan failures must push errors onto error stack**

If a Scan returns false it must push an error onto the scanner's error stack. It is up to the caller to decide to disregard it or not.

* **No regular expressions**

The entire point of creating a recursive descent parser is often to
escape the performance issues and limitations of regular expressions.
For example, regular expressions cannot be inlined by compiler
optimization. Therefore, this package contains *none* of them.

