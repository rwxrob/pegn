package scanner

import "fmt"

var DefaultErrorMessage = `failed to scan`

type Error struct {
	P   int      // can be left blank if Pos is defined
	Pos Position // can be left blank, Report will populate
	Msg string
}

func (e Error) Error() string {
	return fmt.Sprintf("%v at %v", e.Msg, e.Pos)
}

/*
// Error adds an error to the Errors slice. Takes fmt.Sprintf() type
// arguments. The current position (s.Pos) is saved with the error.
// Since s.Pos scans to find the right location if there are multiple
// errors anticipated consider directly appending to Errors instead and
// only using the byte offset position (s.P). Report will detect if the
// first of Errors does not have a Position and will populate them
// efficiently before executing the template. For single errors, calling
// this method should be fine.
func (s *S) Error(a ...any) {
	msg := DefaultErrorMessage
	switch {
	case len(a) > 0:
		msg, _ = a[0].(string)
	case len(a) > 1:
		form, _ := a[0].(string)
		msg = fmt.Sprintf(form, a[1:]...)
	}
	s.Errors = append(s.Errors, Error{Pos: s.Pos(), Msg: msg})
}
*/
