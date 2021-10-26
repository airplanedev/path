# path
Go package for path handling. Currently specialized for, but not restricted to, JSON paths.

## Motivation

This library was written as a way to handle parsing and expressing paths pointing at components in JSON values. Unlike more comprehensive libraries such as [JsonPath](https://github.com/json-path/JsonPath), this library is only intended to handle a single pointer towards one component of a JSON value, rather than being a full query language. On the other hand, it allows for serialization as slices rather than strings, which helps to avoid a class of errors around improperly (de-)escaping path components. 

Each component of a path should either be a string or an integer. In the context of JSON paths, a string represents a value in a JSON object, and an integer represents a value in a JSON array.

## JS Path syntax and examples

When using JS path (de-)serialization, each component of the serialized string should be one of three types:

1. Dot-based strings: These are separated from the previous element by a single `.` (e.g. `a.Test_2009`). They can only consist of alphanumeric characters and the underscore chracter. Note that for the first dot-based string in a serialized path, the dot should be omitted.
2. Square-bracket-and-quote-based strings: These should be contained in square brackets and quotes `[""]` (e.g. `a["Test.ing"]`). These can contain any arbitrary characters, though quotes and backslashes should be backslash-escaped (e.g. `a["\""]`).
3. Integers: These should be contained in square brackets `[]` (e.g. `a[0]`).

```
package main

import (
  "log"

  "github.com/airplane/path"
)

func main() {
  p, err := path.FromJS(`foo[0].bar["baz"]`)

  // Here, p.Components() should be ["foo", 0, "bar", "baz"].
  log.Println(p.Components()) 

  s := p.ToJS()

  // Here, s should be "foo[0].bar.baz". Note that though this is different
  // from the original parsed string, it still will parse back to the same path components.
  log.Println(s) 
}
```

More examples (particularly unusual edge cases) can be found in (javascript_test.go)[https://github.com/airplanedev/path/blob/main/javascript_test.go].
