package lzjson_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-lzjson/lzjson"
)

func TestIsNumJSON(t *testing.T) {
	if want, have := true, lzjson.IsNumJSON([]byte("-1234.56789E+12")); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := true, lzjson.IsNumJSON([]byte("-1234.56789e+12")); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := false, lzjson.IsNumJSON([]byte("-1234.56789A+12")); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestNode_UnmarshalJSON(t *testing.T) {
	str := dummyJSONStr()
	n := lzjson.NewNode()
	var umlr json.Unmarshaler = n
	if err := json.Unmarshal([]byte(str), umlr); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}
	if want, have := str, string(n.Raw()); want != have {
		t.Errorf("\nexpected: %s\ngot: %s", want, have)
	}
}

func TestNode_Unmarshal(t *testing.T) {
	str := dummyJSONStr()
	n, err := lzjson.Decode(strings.NewReader(str))
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}

	type type1 struct {
		Number        float64                `json:"number"`
		String        string                 `json:"string"`
		ArrayOfString []string               `json:"arrayOfString"`
		Object        map[string]interface{} `json:"object"`
	}
	v1 := type1{}

	if err := n.Unmarshal(&v1); err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
		return
	}

	if want, have := 1234.56, v1.Number; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "foo bar", v1.String; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := 4, len(v1.ArrayOfString); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
		return
	}
	if want, have := "one", v1.ArrayOfString[0]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "two", v1.ArrayOfString[1]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "three", v1.ArrayOfString[2]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "four", v1.ArrayOfString[3]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "bar", v1.Object["foo"]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "world", v1.Object["hello"]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := float64(42), v1.Object["answer"]; want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

}
