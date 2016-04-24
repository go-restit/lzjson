package lzjson_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-restit/lzjson"
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

func TestNode_Type(t *testing.T) {

	readJSON := func(str string) (n lzjson.Node) {
		n, err := lzjson.Decode(strings.NewReader(str))
		if err != nil {
			t.Errorf("unexpected error: %#v", err.Error())
			return nil
		}
		return
	}

	if want, have := lzjson.TypeUndefined, (lzjson.NewNode()).Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := lzjson.TypeUndefined, readJSON("").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := lzjson.TypeString, readJSON(`"string"`).Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := lzjson.TypeNumber, readJSON("1234").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := lzjson.TypeNumber, readJSON("-1234.567").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := lzjson.TypeNumber, readJSON("-1234.567E+5").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}

	if want, have := lzjson.TypeObject, readJSON(`{ "foo": "bar" }`).Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := lzjson.TypeArray, readJSON(`[ "foo", "bar" ]`).Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}

	if want, have := lzjson.TypeBool, readJSON("true").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := lzjson.TypeBool, readJSON("false").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
	if want, have := lzjson.TypeNull, readJSON("null").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}

	if want, have := lzjson.TypeError, readJSON("404 not found").Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	}
}

func TestNode_Get(t *testing.T) {
	str := dummyJSONStr()
	n, err := lzjson.Decode(strings.NewReader(str))
	if err != nil {
		t.Errorf("unexpected error: %#v", err.Error())
	}

	if want, have := lzjson.TypeError, n.Get("notExists").Type(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if n := n.Get("number"); n == nil {
		t.Error("unexpected nil value")
	} else if want, have := lzjson.TypeNumber, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := 1234.56, n.Number(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := 1234, n.Int(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := "", n.String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.Bool(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.IsNull(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if n := n.Get("string"); n == nil {
		t.Error("unexpected nil value")
	} else if want, have := lzjson.TypeString, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := float64(0), n.Number(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := "foo bar", n.String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.Bool(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.IsNull(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := 7, n.Len(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	var nilJ lzjson.Node
	if n := n.Get("arrayOfString"); n == nil {
		t.Error("unexpected nil value")
	} else if want, have := lzjson.TypeArray, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := 4, n.Len(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := "one", n.GetN(0).String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := "two", n.GetN(1).String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := "three", n.GetN(2).String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := "four", n.GetN(3).String(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := nilJ, n.GetN(4); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.Bool(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.IsNull(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if n := n.Get("object"); n == nil {
		t.Error("unexpected nil value")
	} else if want, have := lzjson.TypeObject, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if p := n.Get("answer"); p == nil {
		t.Error("unexpected nil value")
	} else if want, have := lzjson.TypeNumber, p.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := false, n.Bool(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.IsNull(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if n := n.Get("true"); n == nil {
		t.Error("unexpected nil value")
	} else if want, have := lzjson.TypeBool, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := true, n.Bool(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if n := n.Get("false"); n == nil {
		t.Error("unexpected nil value")
	} else if want, have := lzjson.TypeBool, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := false, n.Bool(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	} else if want, have := false, n.IsNull(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}

	if n := n.Get("null"); n == nil {
		t.Error("unexpected nil value")
	} else if want, have := lzjson.TypeNull, n.Type(); want != have {
		t.Errorf("expected %s, got %s", want, have)
	} else if want, have := true, n.IsNull(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
