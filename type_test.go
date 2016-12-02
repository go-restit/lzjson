package lzjson_test

import (
	"fmt"
	"testing"

	"github.com/go-restit/lzjson"
)

func TestType(t *testing.T) {
	if want, have := fmt.Sprintf("%#v", lzjson.TypeString.String()), fmt.Sprintf("%#v", lzjson.TypeString.String()); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
