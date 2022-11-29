# Scan functions must advance on success

If a scanner function successfully matches it must leave the scanner pointing to the next byte in the bytes buffer (not yet scanned) ready for the next scan function.
