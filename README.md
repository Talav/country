# country

[![tag](https://img.shields.io/github/tag/talav/country.svg)](https://github.com/talav/country/tags)
[![Go Reference](https://pkg.go.dev/badge/github.com/talav/country.svg)](https://pkg.go.dev/github.com/talav/country)
[![CI](https://github.com/talav/country/workflows/Country%20CI/badge.svg)](https://github.com/talav/country/actions)
[![codecov](https://codecov.io/gh/talav/country/graph/badge.svg)](https://codecov.io/gh/talav/country)
[![License](https://img.shields.io/github/license/talav/country)](./LICENSE)

A country value for Go: an ISO 3166-1 alpha-2 code with its alpha-3 code, numeric code and English name, that cannot hold an invalid value.

## Features

- **ISO 3166-1 built in**: all 249 officially assigned codes, with alpha-3 codes, numeric codes and English short names
- **Named values**: `country.US`, `country.DE`, ... with no error to handle
- **Valid by construction**: an unknown code is an error, never a country
- **Any input form**: alpha-2, alpha-3 or numeric code, in any letter case
- **Comparable**: values work with `==` and as map keys
- **No allocations** when copying, comparing, hashing or parsing
- **Text, JSON and SQL**: `encoding.TextMarshaler`, `driver.Valuer` and `sql.Scanner`, as the alpha-2 code
- **No dependencies** (except testing)

## Installation

```bash
go get github.com/talav/country
```

Requires Go 1.27 or later.

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/talav/country"
)

func main() {
    fmt.Println(country.US.Alpha3())      // USA
    fmt.Println(country.US.NumericCode()) // 840
    fmt.Println(country.US.Name())        // United States of America

    de, err := country.NewCountryFromString("DEU")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(de, de == country.DE) // DE true
}
```

## Usage

### Named values

Every country is a variable named by its alpha-2 code:

```go
country.JP.Alpha2()      // "JP"
country.JP.Alpha3()      // "JPN"
country.JP.NumericCode() // 392
country.JP.Name()        // "Japan"
```

### Parsing

`NewCountryFromString` reads any of the three ISO 3166-1 codes, so data from sources that use different forms gives the same value:

```go
a, _ := country.NewCountryFromString("GB")  // alpha-2
b, _ := country.NewCountryFromString("gbr") // alpha-3, any letter case
c, _ := country.NewCountryFromString("826") // numeric

a == b && b == c && c == country.GB // true
```

White space around the code is ignored. A numeric code has three digits, so Albania is `"008"`, not `"8"`.

Only codes that ISO 3166-1 officially assigns are countries. Reserved, withdrawn and user-assigned codes are errors:

```go
_, err := country.NewCountryFromString("UK")
// invalid country "UK": not an ISO 3166-1 code
```

This also rejects `EU`, `XK` (Kosovo, user-assigned), `AN` (Netherlands Antilles, withdrawn) and UN M49 region codes such as `419`.

### Comparing

A country is a small comparable value:

```go
weights := map[country.Country]float64{
    country.US: 0.6,
    country.JP: 0.4,
}
```

### JSON and SQL

A country is written as its alpha-2 code, in JSON values and map keys alike:

```go
type Fund struct {
    Domicile country.Country         `json:"domicile"`
    Weights  map[country.Country]int `json:"weights"`
}

// {"domicile":"IE","weights":{"JP":40,"US":60}}
```

Reading accepts any code that `NewCountryFromString` accepts. `Scan` accepts a `string` or a `[]byte`. `NULL` is an error; scan a nullable column into a `sql.Null[country.Country]`.

## Zero Value

`country.Country{}` is not a country. It stands for a missing one: `IsValid` returns false, and its codes and name are empty.

## ISO 3166-1 Data

The table in [`iso3166.go`](./iso3166.go) holds the code elements that ISO 3166-1 officially assigns. Names are the English short names in title case, such as `Korea, Republic of`, without the article that ISO writes in parentheses, such as `(the)` in `Bahamas (the)`. A name can change while the codes stay, as Turkey became Türkiye in 2022, so compare countries rather than names.

## Testing

```bash
go test -race -cover ./...
```

## License

MIT License. See [LICENSE](./LICENSE).
