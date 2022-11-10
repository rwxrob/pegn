# PEGN Scanner, Parser, and AST in Go

**🤙Hey there, want to join the PEGN party? 🎉 Now accepting PRs for new Rule implementations. Have a question about how? Open an issue with your question and we'll help you until we can write a contributors guide. Thanks!**

See the [pegn.go](pegn.go) file for documentation of this Go package.

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

* **No regular expressions**

The entire point of creating a recursive descent parser is often to
escape the performance issues and limitations of regular expressions.
For example, regular expressions cannot be inlined by compiler
optimization. Therefore, this package contains *none* of them.
