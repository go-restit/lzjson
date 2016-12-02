//go:generate stringer -type=ParseError -output=error_string.go

package lzjson

// ParseError describe error natures in parsing process
type ParseError int

// types of error
const (
	ErrorUndefined ParseError = iota
	ErrorNotObject
	ErrorNotArray
)

func (err ParseError) Error() string {
	switch err {
	case ErrorUndefined:
		return "undefined"
	case ErrorNotObject:
		return "not an object"
	case ErrorNotArray:
		return "not an array"
	}
	return "unknown parse error"
}

// GoString implements fmt.GoStringer
func (err ParseError) GoString() string {
	return "lzjson." + err.String()
}

// Error is the generic error for parsing
type Error struct {
	Path string
	Err  error
}

// Error implements error type
func (err Error) Error() string {
	if err.Path != "" {
		return err.Path + ": " + err.Err.Error()
	}
	return err.Err.Error()
}

// String implements Stringer
func (err Error) String() string {
	return err.Error()
}
