// Package country provides Country: a country or territory identified by its
// ISO 3166-1 alpha-2 code.
//
// Every code that ISO 3166-1 officially assigns is built in, with its alpha-3
// code, numeric code and English short name. A Country is a named value, such
// as US, or comes from NewCountryFromString.
package country

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

// invalidCodeFormat is the message of the error for a code that ISO 3166-1
// does not assign.
const invalidCodeFormat = "invalid country %q: not an ISO 3166-1 code"

// A Country is a country or territory that ISO 3166-1 officially assigns a
// code to, identified by its alpha-2 code.
//
// A Country is a named value or comes from NewCountryFromString, so it is
// always one that ISO 3166-1 assigns. Countries are comparable with == and
// usable as map keys.
//
// The zero value is not a country. It stands for a missing one: IsValid
// reports false for it, and its codes and name are empty.
type Country struct {
	// code is the ISO 3166-1 alpha-2 code: two upper-case letters. It is
	// empty only in the zero value.
	code string
}

// NewCountryFromString returns the country with the code s, which may be its
// alpha-2 code (US), its alpha-3 code (USA) or its numeric code (840). Letters
// may be in any case, and white space around the code is ignored: "us",
// " USA " and "840" all give US. A numeric code has three digits, so Albania
// is "008".
//
// It returns an error if ISO 3166-1 does not officially assign the code.
// Reserved and user-assigned codes, such as UK, EU and XK, are not countries.
func NewCountryFromString(s string) (Country, error) {
	c, ok := lookup(s)
	if !ok {
		return Country{}, fmt.Errorf(invalidCodeFormat, s)
	}

	return c, nil
}

// IsValid reports whether c is a country. It is false only for the zero
// value, which stands for a missing one.
func (c Country) IsValid() bool {
	return c != Country{}
}

// Alpha2 returns the ISO 3166-1 alpha-2 code of the country, such as US.
func (c Country) Alpha2() string { return c.code }

// Alpha3 returns the ISO 3166-1 alpha-3 code of the country, such as USA.
func (c Country) Alpha3() string { return byCode[c.code].alpha3 }

// NumericCode returns the ISO 3166-1 numeric code of the country, which is
// written with three digits: 840 for US, 008 for AL.
func (c Country) NumericCode() int { return byCode[c.code].numeric }

// Name returns the ISO 3166-1 English short name of the country, such as
// "United States of America", for people to read. Names change, as Turkey
// became Türkiye in 2022, while the codes stay; compare countries, not names.
func (c Country) Name() string { return byCode[c.code].name }

// String returns the alpha-2 code of the country.
func (c Country) String() string { return c.code }

// MarshalText returns the alpha-2 code of the country. It never fails. It
// makes Country an encoding.TextMarshaler, which encoding/json uses for values
// and for map keys.
func (c Country) MarshalText() ([]byte, error) {
	return []byte(c.code), nil
}

// UnmarshalText sets c to the country that NewCountryFromString reads from
// text. If the text is not such a code, it returns the error and leaves c as
// it was. It makes *Country an encoding.TextUnmarshaler.
func (c *Country) UnmarshalText(text []byte) error {
	parsed, ok := lookup(string(text))
	if !ok {
		return fmt.Errorf(invalidCodeFormat, text)
	}

	*c = parsed

	return nil
}

// Value returns the alpha-2 code of the country for use as a query argument.
// It makes Country a driver.Valuer.
func (c Country) Value() (driver.Value, error) {
	return c.code, nil
}

// Scan sets c from a database column value. It makes *Country a sql.Scanner.
//
// A string or []byte must be a code that NewCountryFromString accepts. NULL is
// an error; scan a nullable column into a sql.Null[Country]. On error c is
// unchanged.
func (c *Country) Scan(src any) error {
	switch value := src.(type) {
	case string:
		parsed, err := NewCountryFromString(value)
		if err != nil {
			return err
		}

		*c = parsed

		return nil
	case []byte:
		return c.UnmarshalText(value)
	case nil:
		return errors.New("invalid country: NULL")
	default:
		return fmt.Errorf("invalid country: unsupported type %T", src)
	}
}

// lookup returns the country with the alpha-2, alpha-3 or numeric code s,
// read as NewCountryFromString reads it.
//
// It does not allocate, and s does not escape, so a caller may pass a string
// converted from bytes without a copy on the heap. The code is upper-cased
// into an array, and indexing a map with such an array converted to a string
// does not allocate either. Only ASCII letters change case, so that no other
// letter turns into one a code may hold.
func lookup(s string) (Country, bool) {
	s = strings.TrimSpace(s)

	var key [3]byte
	if len(s) < 2 || len(s) > len(key) {
		return Country{}, false
	}

	for i := range len(s) {
		b := s[i]
		if b >= 'a' && b <= 'z' {
			b = b - 'a' + 'A'
		}

		key[i] = b
	}

	e, ok := byCode[string(key[:len(s)])]

	return Country{code: e.alpha2}, ok
}
