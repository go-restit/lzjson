package lzjson_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-restit/lzjson"
)

func dummyBody() io.Reader {
	return strings.NewReader(`{
    "hello": [
      {
        "name": "world 1",
        "size": 123
      },
      {
        "name": "world 2"
      },
      {
        "name": "world 3"
      }
    ]
  }`)
}

type Namer struct {
	Name string `json:"name"`
}

func Example() {
	body := dummyBody()
	data := lzjson.Decode(body)

	// reading a certain node in the json is straight forward
	fmt.Println(data.Get("hello").GetN(1).Get("name").String())      // output "world 2"
	fmt.Printf("%#v\n", data.Get("hello").GetN(0).Get("size").Int()) // output "123"

	// you may unmarshal the selected child item without defining parent struct
	var namer Namer
	data.Get("hello").GetN(2).Unmarshal(&namer)
	fmt.Println(namer.Name) // output "world 3"

	// you may count elements in an array without defining the array type
	if err, count := data.Get("hello").ParseError(), data.Get("hello").Len(); err == nil {
		fmt.Printf("numbers of item in json.hello: %#v\n", count) // output "numbers of item in json.hello: 3"
	}

	// parse errors inherit along the path, no matter how deep you went
	if err := data.Get("foo").GetN(0).ParseError(); err != nil {
		fmt.Println(err.Error()) // output "json.foo: undefined"
	}
	if err := data.Get("hello").Get("notexists").Get("name").ParseError(); err != nil {
		fmt.Println(err.Error()) // output "json.hello: not an object"
	}
	if err := data.Get("hello").GetN(0).Get("notexists").Get("name").ParseError(); err != nil {
		fmt.Println(err.Error()) // output "json.hello[0].notexists: undefined"
	}
	if err := data.Get("hello").GetN(10).GetN(0).Get("name").ParseError(); err != nil {
		fmt.Println(err.Error()) // output "json.hello[10]: undefined"
	}

	// Output:
	// world 2
	// 123
	// world 3
	// numbers of item in json.hello: 3
	// json.foo: undefined
	// json.hello: not an object
	// json.hello[0].notexists: undefined
	// json.hello[10]: undefined
}
