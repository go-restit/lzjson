package lzjson

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"regexp"
)

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
	TypeUnknown   Type = -1
	TypeUndefined Type = iota
	TypeString
	TypeNumber
	TypeObject
	TypeArray
	TypeBool
	TypeNull
)

// String returns string representations of
// the Type value
func (t Type) String() string {
	switch t {
	case TypeUndefined:
		return "TypeUndefined"
	case TypeString:
		return "TypeString"
	case TypeNumber:
		return "TypeNumber"
	case TypeObject:
		return "TypeObject"
	case TypeArray:
		return "TypeArray"
	case TypeBool:
		return "TypeBool"
	case TypeNull:
		return "TypeNull"
	}
	return "TypeUnknown"
}

// reNumber is the regular expression to match
// any JSON number values
var reNum = regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+\-]?\d+)?$`)

// IsNumJSON test a string and see if it match the
// JSON definition of number
func IsNumJSON(b []byte) bool {
	return reNum.Match(b)
}

// Node is an interface for all JSON nodes
type Node interface {

	// Unmarshal parses the JSON node data into variable v
	Unmarshal(v interface{}) error

	// UnmarshalJSON implements json.Unmarshaler
	UnmarshalJSON(b []byte) error

	// Raw returns the raw JSON string in []byte
	Raw() []byte
}

// NewNode returns an initialized empty Node value
// ready for unmarshaling
func NewNode() Node {
	return &rootNode{}
}

// Decode read and decodes a JSON from io.Reader then
// returns a Node of it
func Decode(reader io.Reader) (n Node, err error) {
	b, err := ioutil.ReadAll(reader)
	n = &rootNode{b}
	return
}

// rootNode is the default implementation of Node
type rootNode struct {
	buf []byte
}

// Unmarshal implements Node
func (n *rootNode) Unmarshal(v interface{}) error {
	return json.Unmarshal(n.buf, v)
}

// UnmarshalJSON implements Node
func (n *rootNode) UnmarshalJSON(b []byte) error {
	n.buf = b
	return nil
}

// Raw implements Node
func (n *rootNode) Raw() []byte {
	return n.buf
}
