# Nonblocking Once

This is a tiny utility for calling a function exactly one time and not blocking
if it's already been called. This utility is useful for actions that you want
to happen exactly once within a certain scope.

## Example

```golang
package main

import (
  "fmt"
  "github.com/bradhe/nbonce"
)

func main() {
  f := func() { fmt.Println("Hello, World!") }
  once := nbonce.NonblockingOnce{}

  for i := 0; i < 10000; i++ {
    once.Once(f)  
  }

  // Wait for the function to finish. If it hasn't been scheduled, returns
  // immediately.
  once.Wait()
}
```
