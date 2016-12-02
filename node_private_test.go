package lzjson

import "testing"

func TestIsNumJSON(t *testing.T) {
	if want, have := true, isNumJSON([]byte("-1234.56789E+12")); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := true, isNumJSON([]byte("-1234.56789e+12")); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
	if want, have := false, isNumJSON([]byte("-1234.56789A+12")); want != have {
		t.Errorf("expected %#v, got %#v", want, have)
	}
}
