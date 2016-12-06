# lzjson [![GoDoc][godoc-badge]][godoc] [![Travis CI results][travis-badge]][travis] [![AppVeyor][appveyor-badge]][appveyor] [![Coverage Status][coveralls-badge]][coveralls]


**lzjson** is a JSON decoding library aims to make you lazy.

Golang default [JSON library](https://golang.org/pkg/encoding/json/) requires to
provide certain data structure (e.g. struct) to decode data to. It is hard to
write type-inspecific code to examine JSON data structure. It is also hard to
determine the abcence or prescence of data field.

This library provide flexible interface for writing generic JSON parsing code.

Key features:

  * zero knowledge parsing: can read and examine JSON structure without
    pre-defining the data structure before hand.

  * lazy parsing: allow you to parse only a specific node into golang data
    structure.

  * compatibility: totally compatible with the default json library

[godoc]: https://godoc.org/github.com/go-restit/lzjson
[godoc-badge]: https://godoc.org/github.com/go-restit/lzjson?status.svg
[travis]: https://travis-ci.org/go-restit/lzjson?branch=master
[travis-badge]: https://api.travis-ci.org/go-restit/lzjson.svg?branch=master
[appveyor]: https://ci.appveyor.com/project/yookoala/lzjson?branch=master
[appveyor-badge]: https://ci.appveyor.com/api/projects/status/github/go-restit/lzjson?branch=master&svg=true
[coveralls]: https://coveralls.io/github/go-restit/lzjson?branch=master
[coveralls-badge]: https://coveralls.io/repos/github/go-restit/lzjson/badge.svg?branch=master

## Example Use

### Decode a JSON

Decode is straight forward with any [io.Reader](https://golang.org/pkg/io/#Reader)
implementation (e.g.
[http.Request.Body](https://golang.org/pkg/net/http/#Request),
[http.Response.Body](https://golang.org/pkg/net/http/#Response),
[strings.Reader](https://golang.org/pkg/strings/#Reader)).

For example, in a [http.HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc):

```go

import (
  "net/http"

  "github.com/go-restit/lzjson"
)


func handler(w http.ResponseWriter, r *http.Request) {
  json := lzjson.Decode(r.Body)
  ...
  ...
}

```

Or as a client:

```go
func main() {
  resp, _ := http.Get("http://foobarapi.com/things")
  json := lzjson.Decode(resp.Body)
  ...
  ...
}
```

### Get a node in an object or an array

You may retrieve the JSON value of any node.

```go
// get "foo" in the json
foo := json.Get("foo")

// get the 10th item in foo
// (like ordinary array, 0 is the first)
item10 := foo.GetN(9)
```

### Every node knows what it is

```go
body := strings.NewReader(`
{
  "string": "hello world",
  "number": 3.14,
  "bool": true,
  "array": [1, 2, 3, 5],
  "object": {"foo": "bar"}
}
`)
json := lzjson.Decode(body)

fmt.Printf("%s", json.Get("string").Type()) // output "TypeString"
fmt.Printf("%s", json.Get("number").Type()) // output "TypeNumber"
fmt.Printf("%s", json.Get("bool").Type())   // output "TypeBool"
fmt.Printf("%s", json.Get("array").Type())  // output "TypeArray"
fmt.Printf("%s", json.Get("object").Type()) // output "TypeObject"
```

### Evaluating values a JSON node

For basic value types (string, int, bool), you may evaluate them directly.

```go
code := json.Get("code").Int()
message := json.Get("message").String()
```

### Partial Unmarsaling

You may decode only a child-node in a JSON structure.

```go

type Item struct {
  Name   string `json:"name"`
  Weight int    `json:"weight"`
}

var item Item
item10 := foo.GetN(9)
item10.Unmarshal(&item)
log.Printf("item: name=%s, weight=%d", item.Name, item.Weight)

```

### Chaining

You may chain `Get` and `GetN` to get somthing deep within.

```go

helloIn10thBar := lzjson.Decode(r.Body).Get("foo").GetN(9).Get("hello")

```

### Looping Object or Array

Looping is straight forward with `Len` and `GetKeys`.

```go
var item Item
for i := 0; i<foo.Len(); i++ {
  foo.Get(i).Unmarshal(&item)
  log.Printf("i=%d, value=%#v", i, item)
}

for _, key := range json.GetKeys() {
  log.Printf("key=%#v, value=%#v", key, json.Get(key).String())
}
```

### Error knows their location

With chaining, it is important where exactly did any parse error happen.

```go

body := strings.NewReader(`
{
  "hello": [
    {
      "name": "world 1"
    },
    {
      "name": "world 2"
    },
    {
      "name": "world 3"
    },
  ],
}
`)
json := lzjson.Decode(body)

inner := json.Get("hello").GetN(2).Get("foo").Get("bar").GetN(0)
if err := inner.ParseError(); err != nil {
  fmt.Println(err.Error()) // output: "hello[2].foo: undefined"
}

```

### Full Example

Put everything above together, we can do something like this:

```go

package main

import (
  "log"
  "net/http"

  "github.com/go-restit/lzjson"
)

type Thing struct {
  ID        string    `json:"id"`
  Name      string    `json:"name"`
  Found     time.Time `json:"found"`
  FromEarth bool      `json:"from_earth"`
}

/**
 * assume the API endpoints returns data:
 * {
 *   "code": 200,
 *   "data": [
 *     ...
 *   ]
 * }
 *
 * or error:
 * {
 *   "code": 500,
 *   "message": "some error message"
 * }
 *
 */
func main() {
  resp, err := http.Get("http://foobarapi.com/things")
  if err != nil {
    panic(err)
  }

  // decode the json as usual, if no error
  json := lzjson.Decode(resp.Body)
  if code := json.Get("code").Int(); code != 200 {
    message := json.Get("message").String()
    log.Fatalf("error %d: ", code, message)
  }

  // get the things array
  things := json.Get("data")

  // loop through the array
  for i := 0; i<things.Len(); i++ {
    thing := things.GetN(i)
    if err := thing.ParseError(); err != nil {
      log.Fatal(err.Error())
    }

    // if the thing is not from earth, unmarshal
    // as a struct then read the details
    if !thing.Get("from_earth").Bool() {
      var theThing Thing
      thing.Unmarshal(&theThing)
      log.Printf("Alien found! %#v", theThing)
    }
  }

}

```

For more details, please read the [documentation][godoc]


## Contirbuting

Your are welcome to contribute to this library.

To report bug, please use the [issue tracker][issue tracker].

To fix an existing bug or implement a new feature, please:

1. Check the [issue tracker][issue tracker] and [pull requests][pull requests] for existing discussion.
2. If not, please open a new issue for discussion.
3. Write tests.
4. Open a pull request referencing the issue.
5. Have fun :-)

[issue tracker]: https://github.com/go-restit/lzjson/issues
[pull requests]: https://github.com/go-restit/lzjson/pulls


## Licence

This software is licenced with the [MIT Licence] [licence].
You can obtain a copy of the licence in this repository.

[licence]: /LICENCE
