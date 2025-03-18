# Layered Map

It is an algorithm that organizes data by its key, like some map structures.

The differential with `Layered Map` is that the time used to search for some key is equivalent, when not allocating more memory, to the size in bytes of the key.

It means that if you have one key called "four", it will take 4 times to access the value stored in "four".

It supports everything to be a key because the key is one generic slice of bytes, so you can use `strings`, `emojis 😊`, `json` and `media` files if you want, just converting it to slice of bytes.

# How it works

For each byte in the given key, the algorithm will create a way with linked lists to give access to the position or value in the last byte of the key, following the sequence of the key provided. The stored value can be anything, the user should know what is the interface that will handle the bytes.
Ex: I want to store "Hello, World!" with the key "abc"

```
[a] [ ] [ ]
 |___
[ ] [b] [ ]
     |___
[ ] [ ] [c] -> stores "Hello, World!"
```

Explaining: For each byte, if the next layer of chars isn't allocated, the algorithm create the next layer and try to store for the next, making it recursive while it dont reach the end of the key.

# Pros and Cons

## Pros

The `Layered Map` is particularly useful when dealing with large datasets, such as document-oriented data with millions of entries. Unlike traditional search methods that require scanning the entire database to find a matching document—which can be time-consuming—this algorithm retrieves data in a time proportional to the size of the provided key (in bytes). This makes it highly efficient for scenarios where quick lookups based on specific keys are critical.

## Cons

On the downside, if your dataset contains a large number of unique keys, the memory usage can increase significantly. Each distinct key creates its own path of nodes in the structure, leading to higher memory allocation. This trade-off may become a concern in applications with highly diverse keys or limited memory resources.

# Usage Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/bethecozmo/layeredmap"
)

func main() {
  lm := layeredmap.New()
  ttl := time.Second * 5

  // Add user Alice with expiration time
  lm.Add([]byte("user"), "Alice", &ttl)

  // Add user Robert with no expiration time
  lm.Add([]byte("user"), "Robert", nil)

  // Get all stored values with key `user`
  values, found := lm.GetAll([]byte("user"))
  if found {
    fmt.Println("Values:", values) // [Alice Roberto]
  } else {
    fmt.Println("Key not found or all values have expired")
  }

  // Remove the last stored value and return it
  last, found := lm.PopLast([]byte("user"))
  if found {
    fmt.Println("Last value:", last) // Roberto
  }
}
```

See more on [layeredmap.go](./layeredmap.go)