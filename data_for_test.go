package lzjson_test

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandString generate fix length random strings
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// dummyJSONStr returns dummy JSON string for test
func dummyJSONStr() string {
	return `{
    "number": 1234.56,
    "string": "foo bar",
    "arrayOfString": [
      "one",
      "two",
      "three",
      "four"
    ],
    "object": {
      "foo": "bar",
      "hello": "world",
      "answer": 42
    },
    "true": true,
    "false": false,
    "null": null
  }`
}
