# Design considerations

Here's the thinking behind the design decisions as they were made.

* [Use `C_` prefix for classes to preserve identifiers](../3?L)
* [Simplest rune pegn.Scanner interface](../4?L)
* [Decouple meta-data association](../5?L)
* [All scanner functions advance the scanner themselves](../6?L)
* [Scan functions must not change state on failure](../7?L)
* [Scan functions must advance on success](../8?L)
* [Scan failures must push errors onto error stack](../9?L)
* [No regular expressions](../10?L)
