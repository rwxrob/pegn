# Scan failures must push errors onto error stack

If a Scan returns false it must push an error onto the scanner’s error stack. It is up to the       caller to decide to disregard it or not.
