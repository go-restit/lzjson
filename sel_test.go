package lzjson

import "testing"

func TestSelLex(t *testing.T) {
	type testPair struct {
		Sel      string
		Expected []selItem
	}

	tests := []testPair{
		testPair{
			"hello",
			[]selItem{
				selItem{typ: selItemProp, val: "hello"},
			},
		},
		testPair{
			"hello.world",
			[]selItem{
				selItem{typ: selItemProp, val: "hello"},
				selItem{typ: selItemDot, val: "."},
				selItem{typ: selItemProp, val: "world"},
			},
		},
		testPair{
			"hello[123]world",
			[]selItem{
				selItem{typ: selItemProp, val: "hello"},
				selItem{typ: selItemLeftBrac, val: "["},
				selItem{typ: selItemNumber, val: "123"},
				selItem{typ: selItemRightBrac, val: "]"},
				selItem{typ: selItemProp, val: "world"},
			},
		},
		testPair{
			"hello.world[123]",
			[]selItem{
				selItem{typ: selItemProp, val: "hello"},
				selItem{typ: selItemDot, val: "."},
				selItem{typ: selItemProp, val: "world"},
				selItem{typ: selItemLeftBrac, val: "["},
				selItem{typ: selItemNumber, val: "123"},
				selItem{typ: selItemRightBrac, val: "]"},
			},
		},
		testPair{
			"hello[12][3]world.foo.bar[4]",
			[]selItem{
				selItem{typ: selItemProp, val: "hello"},
				selItem{typ: selItemLeftBrac, val: "["},
				selItem{typ: selItemNumber, val: "12"},
				selItem{typ: selItemRightBrac, val: "]"},
				selItem{typ: selItemLeftBrac, val: "["},
				selItem{typ: selItemNumber, val: "3"},
				selItem{typ: selItemRightBrac, val: "]"},
				selItem{typ: selItemProp, val: "world"},
				selItem{typ: selItemDot, val: "."},
				selItem{typ: selItemProp, val: "foo"},
				selItem{typ: selItemDot, val: "."},
				selItem{typ: selItemProp, val: "bar"},
				selItem{typ: selItemLeftBrac, val: "["},
				selItem{typ: selItemNumber, val: "4"},
				selItem{typ: selItemRightBrac, val: "]"},
			},
		},
		testPair{
			"hello.world[\"foo\"][\"bar\"]",
			[]selItem{
				selItem{typ: selItemProp, val: "hello"},
				selItem{typ: selItemDot, val: "."},
				selItem{typ: selItemProp, val: "world"},
				selItem{typ: selItemLeftBrac, val: "["},
				selItem{typ: selItemString, val: "foo"},
				selItem{typ: selItemRightBrac, val: "]"},
				selItem{typ: selItemLeftBrac, val: "["},
				selItem{typ: selItemString, val: "bar"},
				selItem{typ: selItemRightBrac, val: "]"},
			},
		},
		testPair{
			"[\"foo\"][\"bar\"]",
			[]selItem{
				selItem{typ: selItemLeftBrac, val: "["},
				selItem{typ: selItemString, val: "foo"},
				selItem{typ: selItemRightBrac, val: "]"},
				selItem{typ: selItemLeftBrac, val: "["},
				selItem{typ: selItemString, val: "bar"},
				selItem{typ: selItemRightBrac, val: "]"},
			},
		},
		testPair{
			"[\"foo and \\\"bar\\\"\"]",
			[]selItem{
				selItem{typ: selItemLeftBrac, val: "["},
				selItem{typ: selItemString, val: "foo and \\\"bar\\\""},
				selItem{typ: selItemRightBrac, val: "]"},
			},
		},
		testPair{
			"['foo']['bar']",
			[]selItem{
				selItem{typ: selItemLeftBrac, val: "["},
				selItem{typ: selItemString, val: "foo"},
				selItem{typ: selItemRightBrac, val: "]"},
				selItem{typ: selItemLeftBrac, val: "["},
				selItem{typ: selItemString, val: "bar"},
				selItem{typ: selItemRightBrac, val: "]"},
			},
		},
		testPair{
			"['foo\\'s bar']",
			[]selItem{
				selItem{typ: selItemLeftBrac, val: "["},
				selItem{typ: selItemString, val: "foo\\'s bar"},
				selItem{typ: selItemRightBrac, val: "]"},
			},
		},
	}

	for _, test := range tests {
		i, l := 0, len(test.Expected)
		lex := lexSel(test.Sel)
		go lex.run()

		for v := lex.nextItem(); v.typ != selItemEnd; v = lex.nextItem() {
			if i >= l {
				t.Errorf("sel=%#v pos=%#v error=\"index out of range\" got=%#v", test.Sel, i, v)
			} else if want, have := test.Expected[i], v; want != have {
				t.Errorf("sel=%#v pos=%#v expected=%#v got=%#v", test.Sel, i, want, have)
			}
			i++
		}
		if want, have := l, i; want != have {
			t.Errorf("sel=%#v number of output mismatch. expected %#v, got %#v", test.Sel, want, have)
		}
	}

}
