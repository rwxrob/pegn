# Simple pegn.Scanner interface

the temptation to add more methods to the pegn.Scanner interface has been constant, and overcome. The driving design principle for it is that it could be used entirely without any other dependency on PEGN for *any* rune scanning. This is why ScannerErrors interface methods use simple error and not ScanError (even though ScanError is strongly recommended with implementing pegn.Scanners).
