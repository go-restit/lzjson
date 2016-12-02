package lzjson_test

import (
	"fmt"
	"testing"

	"github.com/go-restit/lzjson"
)

func TestParseError_GoString(t *testing.T) {
	if want, have := "lzjson.ErrorUndefined", fmt.Sprintf("%#v", lzjson.ErrorUndefined); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "lzjson.ErrorNotObject", fmt.Sprintf("%#v", lzjson.ErrorNotObject); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "lzjson.ErrorNotArray", fmt.Sprintf("%#v", lzjson.ErrorNotArray); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "unknown parse error", lzjson.ParseError(-1).Error(); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := "lzjson.ParseError(-1)", fmt.Sprintf("%#v", lzjson.ParseError(-1)); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}

func TestError_String(t *testing.T) {
	err := lzjson.Error{
		Path: "hello",
		Err:  fmt.Errorf("some error msg"),
	}
	if want, have := "hello: some error msg", err.String(); want != have {
		t.Errorf("\nexpected:\n%s\ngot:\n%s", want, have)
	}
}
