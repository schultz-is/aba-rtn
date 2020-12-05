# rtnutil

This package provides utilities for working with American Bankers Association
routing transit numbers; also known as ABA RTNs.

More information about the composition, structure, and usage of ABA RTNs can be
found on the relevant Wikipedia page[^1].

[^1]: https://en.wikipedia.org/wiki/ABA_routing_transit_number

## Installation

To install as a dependency in a go project:

```console
go get github.com/schultz-is/go-threefish
```

## Usage

### Validating an RTN

The format and checksum of an RTN can be validated via the `Validate`
package-level function. Note that a successful result from this function does
not necessarily mean that the RTN is assigned and in use!

```go
package main

import (
  "fmt"

  "github.com/schultz-is/rtnutil"
)

func main() {
  err := rtnutil.Validate("044000037")
  if err != nil {
    panic(err)
  }

  fmt.Println("RTN is valid!")
}
```

### Calculating a missing RTN digit

In the case where an RTN is missing a check digit or one of the digits is
illegible, it is possible to calculate the missing digit by replacing the
missing digit with an "X" and passing it to the `GetMissingDigit` function.
Only a single missing digit can be calculated this way, and the RTN provided
must be missing exactly one digit.

```go
package main

import (
  "fmt"

  "github.com/schultz-is/rtnutil"
)

func main() {
  missingDigit, err := rtnutil.GetMissingDigit("04400003X")
  if err != nil {
    panic(err)
  }

  fmt.Printf("the missing digit is %d", missingDigit)
}
```

## Testing

Unit tests can be run and test coverage can be viewed via the provided
Makefile.

```console
make test
make coverage
```
