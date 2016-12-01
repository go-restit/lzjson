//go:generate stringer -type=Type

package lzjson

// Type represents the different type of JSON values
// (string, number, object, array, true, false, null)
// true and false are combined as bool for obvious reason
type Type int

// These constant represents different JSON value types
// as specified in http://www.json.org/
// with some exception:
// 1. true and false are combined as bool for obvious reason; and
// 2. TypeUnknown for empty strings
const (
	TypeError     Type = -1
	TypeUndefined Type = iota
	TypeString
	TypeNumber
	TypeObject
	TypeArray
	TypeBool
	TypeNull
)

// GoString implements fmt.GoStringer
func (t Type) GoString() string {
	return "lzjson." + t.String()
}
