# PEGN Scanner, Parser, and AST in Go

🚧 Still a long way from 2023-01 version of `pegn` package to accompany the next release of the [PEGN specification](https://github.com/rwxrob/pegn-spec).

See the [pegn.go](pegn.go) file for documentation of this Go package.

## More than just a grammar spec

This repo primarily contains the specification of the PEGN meta language with all the expected parts of a package used to handle a grammar (which some may generate directly from a PEGN file):

* `rule` - rule definitions and meta data
* `rule/id` - rule unique identifiers (neg int for all PEGN)
* `scan` - scan functions (for all rules)
* `parse` - parse functions (that usually call scan functions)
* `tk` - tokens (rune and string constants)
* `is` - class functions (called from scan functions)

This repo also contains some "batteries included" tools to help work with PEGN

* `types.go`  - type constants
* `docs/`     - multi-lingual rule documentation
* `scanner/`  - pegn.Scanner implementation (DefaultScanner)
* `qstack/` - queue and stack data structure combined
* `curs/`     - curs.R struct for position within bytes buffer
* `cmd.go`    - Bonzai stateful command tree composition
* `cmd/pegn`  - helper utility with code generation

## Design Considerations

* **Use `C_` prefix for classes to preserve identifiers**

The use of aliases just wasn't working. This makes them immediately identifiable.

* **Individual scan, parse, and is (class) function files**

By keeping the function types of PEGN in their own files it becomes trivial to create code generators by simply copying the files needed into an isolated package directory and adjusting the package line.

* **Simplest rune pegn.Scanner interface**

The temptation to add more methods to the pegn.Scanner interface has been constant, and overcome. The driving design principle for it is that it could be used entirely without any other dependency on PEGN for *any* rune scanning. This is why ScannerErrors interface methods use simple error and not ScanError (even though ScanError is strongly recommended with implementing pegn.Scanners).

* **Decouple meta-data association**

For a time, associating the PEGN notation, aliases, and human-friendly descriptions of each pegn.Type was done by forcing them all into a pegn.Rule interface implementation, the idea being that such information would travel with the thing doing the scanning and parsing.

However, it became clear that only a reserved pegn.Type is needed in order to accomplish the same thing allowing simple, decoupled ScanFunc and ParseFunc instances rather than a full struct implementation. This allows the help and other documentation to be contained separately and greatly facilitated code generation. (It's worth noting that this was the original 2018 design as well.)

* **All scanner functions advance the scanner themselves**

Since it is not always clear how much of the bytes buffer will be
needed, the task of advancing the scanner (with `s.Scan()` or however)
is up to each scan function.

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

