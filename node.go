package lzjson

import (
	"encoding/json"
	"fmt"
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
	TypeError     Type = -1
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
	return "TypeError"
}

// GoString implements fmt.GoStringer
func (t Type) GoString() string {
	return "lzjson." + t.String()
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

	// Type returns the Type of the containing JSON value
	Type() Type

	// GetKeys gets an array object's keys,
	// or nil if not an object
	GetKeys() []string

	// Get gets object's inner value.
	// Only works with Object value type
	Get(key string) (inner Node)

	// Len gets the length of the value
	// Only works with Array and String value type
	Len() int

	// GetN gets array's inner value.
	// Only works with Array value type.
	// 0 for the first item.
	GetN(nth int) Node

	// String unmarshal the JSON into string then return
	String() (v string)

	// Number unmarshal the JSON into float64 then return
	Number() (v float64)

	// Int unmarshal the JSON into int the return
	Int() (v int)

	// Bool unmarshal the JSON into bool then return
	Bool() (v bool)

	// IsNull tells if the JSON value is null or not
	IsNull() bool

	// Error returns the JSON parse error, if any
	Error() error
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
	n = &rootNode{buf: b}
	return
}

// rootNode is the default implementation of Node
type rootNode struct {
	buf    []byte
	mapBuf map[string]rootNode
	err    error
}

// Unmarshal implements Node
func (n *rootNode) Unmarshal(v interface{}) error {
	return json.Unmarshal(n.buf, v)
}

// UnmarshalJSON implements Node
func (n *rootNode) UnmarshalJSON(b []byte) error {
	n.buf = b
	n.mapBuf = nil
	return nil
}

// Raw implements Node
func (n *rootNode) Raw() []byte {
	return n.buf
}

// Type implements Node
func (n rootNode) Type() Type {

	switch {
	case n.err != nil:
		// for error, return TypeError
		return TypeError
	case n.buf == nil || len(n.buf) == 0:
		// for nil raw, return TypeUndefined
		return TypeUndefined
	case n.buf[0] == '"':
		// simply examine the first character
		// to determine the value type
		return TypeString
	case n.buf[0] == '{':
		// simply examine the first character
		// to determine the value type
		return TypeObject
	case n.buf[0] == '[':
		// simply examine the first character
		// to determine the value type
		return TypeArray
	case string(n.buf) == "true":
		fallthrough
	case string(n.buf) == "false":
		return TypeBool
	case string(n.buf) == "null":
		return TypeNull
	case IsNumJSON(n.buf):
		return TypeNumber
	}

	// return TypeUnknown for all other cases
	return TypeError
}

func (n *rootNode) genMapBuf() error {
	if n.Type() != TypeObject {
		return fmt.Errorf("the node is not an object")
	}
	if n.mapBuf != nil {
		return nil // previously done, use the previous result
	}

	// generate the map
	n.mapBuf = map[string]rootNode{}
	return n.Unmarshal(&n.mapBuf)
}

// GetKeys get object keys of the node.
// If the node is not an object, returns nil
func (n *rootNode) GetKeys() (keys []string) {
	if err := n.genMapBuf(); err != nil {
		return
	}
	keys = make([]string, 0, len(n.mapBuf))
	for key := range n.mapBuf {
		keys = append(keys, key)
	}
	return
}

// Get implements Node
func (n *rootNode) Get(key string) (inner Node) {
	if err := n.genMapBuf(); err != nil {
		inner = &rootNode{err: err} // dump the error
	} else if val, ok := n.mapBuf[key]; !ok {
		inner = &rootNode{err: fmt.Errorf("field %#v not found", key)}
	} else {
		inner = &val
	}
	return
}

// Len gets the length of the value
// Only works with Array and String value type
func (n *rootNode) Len() int {
	switch n.Type() {
	case TypeString:
		return len(string(n.buf)) - 2 // subtact the 2 " marks
	case TypeArray:
		vslice := []*rootNode{}
		n.Unmarshal(&vslice)
		return len(vslice)
	}
	// default return -1 (for type mismatch)
	return -1
}

// GetN implements Node
func (n *rootNode) GetN(nth int) Node {
	if n.Type() != TypeArray {
		return nil
	}

	vslice := []rootNode{}
	n.Unmarshal(&vslice)
	if nth < len(vslice) {
		return &vslice[nth]
	}
	return nil
}

// String implements Node
func (n *rootNode) String() (v string) {
	n.Unmarshal(&v)
	return
}

// Number implements Node
func (n *rootNode) Number() (v float64) {
	n.Unmarshal(&v)
	return
}

// Int implements Node
func (n *rootNode) Int() (v int) {
	return int(n.Number())
}

// Bool implements Node
func (n *rootNode) Bool() (v bool) {
	n.Unmarshal(&v)
	return
}

// IsNull implements Node
func (n *rootNode) IsNull() bool {
	return n.Type() == TypeNull
}

// Error implements Node
func (n *rootNode) Error() error {
	return n.err
}
