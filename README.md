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

  json := lzjson.Decode(resp.Body)
  if code := json.Get("code").Int(); code != 200 {
    message := json.Get("message").String()
    panic(message)
  }

  things := json.Get("data")
  for i := 0; i<=things.Len(); i++ {
    thing := things.GetN(0)
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
