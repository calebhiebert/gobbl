# LUIS Middleware

## What it does

Queries a LUIS endpoing with the contents of `c.Request.Text` and sets some flags.

## Flags Set

- `intent` - A string representing the intent returned from LUIS
- `luis` - The entire result object that LUIS returned. Type `*luis.Response`
- `luis:e:*` - For each entity returned, a flag of []string will be added. For example, if the entity builtin.number is returned, the flag would be `luis:e:builtin.number` -> `[]string{"1"}`

## Usage

1. Create a LUIS instance

``` go
louie, err := luis.New("https://example.luis/endpoint")
if err != nil {
  panic(err)
}
```

2. Use the middleware

```go
bot.Use(luis.Middleware(louie))
```

### Full Example

```go
package main

import(
  "github.com/calebhiebert/gobbl"
  "github.com/calebhiebert/gobbl/luis"
)

func main() {
  gobblr := gbl.New()

  louie, err := luis.New("https://example.luis/endpoint")
  if err != nil {
    panic(err)
  }

  gobblr.Use(luis.Middleware(louie))
}
```