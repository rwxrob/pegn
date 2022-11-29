# Scan functions must not change state on failure

The scanner performing the scan must be in exactly the same state (rune, byte pointer, previous     byte pointer, etc.) if the scan fails (returning `false`). The only exception is pushing one or     more errors into the scanner's memory.
