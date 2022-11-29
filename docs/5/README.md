# Decouple meta data

For a time, associating the PEGN notation, aliases, and human-friendly descriptions of each pegn.Type was done by forcing them all into a pegn.Rule interface implementation, the idea being that such information would travel with the thing doing the scanning and parsing.

However, it became clear that only a reserved pegn.Type is needed in order to accomplish the same thing allowing simple, decoupled ScanFunc and ParseFunc instances rather than a full struct implementation. This allows the help and other documentation to be contained separately and greatly facilitated code generation. (It's worth noting that this was the original 2018 design as well.)

