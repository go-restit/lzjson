package lzjson

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

// reNumber is the regular expression to match
// any JSON number values
var reNum = regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+\-]?\d+)?$`)

// isNumJSON test a string and see if it match the
// JSON definition of number
func isNumJSON(b []byte) bool {
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

	// ParseError returns the JSON parse error, if any
	ParseError() error
}

// NewNode returns an initialized empty Node value
// ready for unmarshaling
func NewNode() Node {
	return &rootNode{}
}

// Decode read and decodes a JSON from io.Reader then
// returns a Node of it
func Decode(reader io.Reader) Node {
	b, err := ioutil.ReadAll(reader)
	return &rootNode{
		buf: b,
		err: err,
	}
}

// rootNode is the default implementation of Node
type rootNode struct {
	path   string
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
	case isNumJSON(n.buf):
		return TypeNumber
	}

	// return TypeUnknown for all other cases
	return TypeError
}

func (n *rootNode) genMapBuf() error {
	if n.Type() != TypeObject {
		return ErrorNotObject
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

func (n *rootNode) keyPath(key string) string {
	fmtKey := "." + key
	if strings.IndexAny(key, " /-") >= 0 {
		fmtKey = fmt.Sprintf("[%#v]", key)
	}
	return n.path + fmtKey
}

// Get implements Node
func (n *rootNode) Get(key string) (inner Node) {

	// the child key path
	path := n.keyPath(key)

	// if there is previous error, inherit
	if err := n.ParseError(); err != nil {
		return &rootNode{
			path: path,
			err:  err,
		}
	}

	if err := n.genMapBuf(); err != nil {
		if err == ErrorNotObject {
			path = n.path // fallback to the parent entity
		}
		inner = &rootNode{
			path: path,
			err: Error{
				Path: "json" + path,
				Err:  err,
			},
		}
	} else if val, ok := n.mapBuf[key]; !ok {
		inner = &rootNode{
			path: path,
			err: Error{
				Path: "json" + path,
				Err:  ErrorUndefined,
			},
		}
	} else {
		val.path = path
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

func (n *rootNode) nthPath(nth int) string {
	return fmt.Sprintf("%s[%d]", n.path, nth)
}

// GetN implements Node
func (n *rootNode) GetN(nth int) Node {

	// the path to nth node
	path := n.nthPath(nth)

	// if there is previous error, inherit
	if err := n.ParseError(); err != nil {
		return &rootNode{
			path: path,
			err:  err,
		}
	}

	if n.Type() != TypeArray {
		return &rootNode{
			path: n.path,
			err: Error{
				Path: "json" + n.path,
				Err:  ErrorNotArray,
			},
		}
	}

	vslice := []rootNode{}
	n.Unmarshal(&vslice)
	if nth < len(vslice) {
		vslice[nth].path = path
		return &vslice[nth]
	}
	return &rootNode{
		path: path,
		err: Error{
			Path: "json" + path,
			Err:  ErrorUndefined,
		},
	}
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

// ParseError implements Node
func (n *rootNode) ParseError() error {
	return n.err
}
