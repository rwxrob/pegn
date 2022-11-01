# PEGN Scanner, Parser, and AST in Go

See the [doc.go](doc.go) for documentation of this Go package.

## Design Considerations

* **All scanner functions advance the scanner themselves**

Since it is not always clear how much of the bytes buffer will be
needed, the task of advancing the scanner (with `s.Scan()` or however)
is up to each scan function.

* **Scan functions must not change state on fail**

The scanner performing the scan must be in exactly the same state (rune,
byte pointer, previous byte pointer, etc.) if the scan fails (returning
`false`).

* **Scan functions must advance on success**

If a scanner function successfully matches it must leave the scanner
pointing to the next byte in the bytes buffer (not yet scanned) ready
for the next scan function.

* **PEGN class names as all caps (like tokens)**

With the requirement to provide an initial cap to export a function in
Go otherwise all lower-case PEGN class names have been forced to all
upper-case. This works out since there is no name collisions between
PEGN classes, tokens, and functions. The documentation disambiguates
which is which.

```go
scan.WS(s)
scan.SPACE(s)
scan.SP(s)
```

* **Add `Some` for greedy scans**

Even though PEG itself calls for all repeating operators to be greedy by
default, sometimes such is not wanted in this package library of scan
functions. Therefore, anything prefixed with `Some` will continue to
scan for as long as the content matches. For example, `SomeWS(s)` will
consume all the white-space it finds from the current position, while
`WS(s)` will only return true if the next rune is white space.

* **No regular expressions**

The entire point of creating a recursive descent parser is often to
escape the performance issues and limitations of regular expressions.
For example, regular expressions cannot be inlined by compiler
optimization. Therefore, this package contains *none* of them.
