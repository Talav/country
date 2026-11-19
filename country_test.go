package country_test

import (
	"encoding"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/talav/country"
)

// parse builds a country from a code that the test knows to be valid.
func parse(t *testing.T, code string) country.Country {
	t.Helper()

	c, err := country.NewCountryFromString(code)
	require.NoError(t, err)

	return c
}

func TestNewCountryFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		code    string
		want    country.Country
		wantErr bool
	}{
		{name: "alpha-2", code: "US", want: country.US},
		{name: "alpha-3", code: "USA", want: country.US},
		{name: "numeric", code: "840", want: country.US},
		{name: "alpha-2 in lower case", code: "us", want: country.US},
		{name: "alpha-3 in mixed case", code: "Deu", want: country.DE},
		{name: "alpha-3 unlike the English name", code: "CHE", want: country.CH},
		{name: "numeric below 100", code: "008", want: country.AL},
		{name: "numeric below 10", code: "004", want: country.AF},
		{name: "surrounding white space", code: " gb\t", want: country.GB},
		{name: "padded alpha-3", code: "\nJPN ", want: country.JP},
		{name: "padded numeric", code: " 826 ", want: country.GB},
		{name: "surrounding non-breaking space", code: " FR ", want: country.FR},
		{name: "first in the table", code: "AD", want: country.AD},
		{name: "last in the table", code: "ZW", want: country.ZW},
		{name: "newest", code: "SSD", want: country.SS},
		{name: "empty", code: "", wantErr: true},
		{name: "only white space", code: "  ", wantErr: true},
		{name: "one letter", code: "U", wantErr: true},
		{name: "four letters", code: "USAA", wantErr: true},
		{name: "space inside", code: "U S", wantErr: true},
		{name: "subdivision", code: "US-CA", wantErr: true},
		{name: "exceptionally reserved", code: "UK", wantErr: true},
		{name: "European Union", code: "EU", wantErr: true},
		{name: "user-assigned", code: "XK", wantErr: true},
		{name: "withdrawn alpha-2", code: "AN", wantErr: true},
		{name: "withdrawn alpha-3", code: "SUN", wantErr: true},
		{name: "unassigned", code: "QQ", wantErr: true},
		{name: "numeric without leading zeros", code: "8", wantErr: true},
		{name: "numeric with too few digits", code: "08", wantErr: true},
		{name: "numeric with four digits", code: "0840", wantErr: true},
		{name: "unassigned numeric", code: "999", wantErr: true},
		{name: "UN M49 region", code: "419", wantErr: true},
		{name: "letter that upper-cases to an ASCII one", code: "uſ", wantErr: true},
		{name: "non-ASCII letter", code: "ÜS", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := country.NewCountryFromString(tt.code)
			if tt.wantErr {
				require.EqualError(t, err, fmt.Sprintf("invalid country %q: not an ISO 3166-1 code", tt.code))
				assert.Equal(t, country.Country{}, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCountryCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		country     country.Country
		wantAlpha2  string
		wantAlpha3  string
		wantNumeric int
		wantName    string
	}{
		{name: "named value", country: country.US, wantAlpha2: "US", wantAlpha3: "USA", wantNumeric: 840, wantName: "United States of America"},
		{name: "parsed", country: parse(t, "usa"), wantAlpha2: "US", wantAlpha3: "USA", wantNumeric: 840, wantName: "United States of America"},
		{name: "numeric below 100", country: country.AL, wantAlpha2: "AL", wantAlpha3: "ALB", wantNumeric: 8, wantName: "Albania"},
		{name: "alpha-3 unlike the English name", country: country.CH, wantAlpha2: "CH", wantAlpha3: "CHE", wantNumeric: 756, wantName: "Switzerland"},
		{name: "non-ASCII name", country: country.CI, wantAlpha2: "CI", wantAlpha3: "CIV", wantNumeric: 384, wantName: "Côte d'Ivoire"},
		{name: "long name", country: country.GB, wantAlpha2: "GB", wantAlpha3: "GBR", wantNumeric: 826, wantName: "United Kingdom of Great Britain and Northern Ireland"},
		{name: "zero value", country: country.Country{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.wantAlpha2, tt.country.Alpha2())
			assert.Equal(t, tt.wantAlpha3, tt.country.Alpha3())
			assert.Equal(t, tt.wantNumeric, tt.country.NumericCode())
			assert.Equal(t, tt.wantName, tt.country.Name())
		})
	}
}

func TestCountryIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		country country.Country
		want    bool
	}{
		{name: "zero value", country: country.Country{}, want: false},
		{name: "named value", country: country.US, want: true},
		{name: "parsed", country: parse(t, "JPN"), want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.country.IsValid())
		})
	}
}

func TestCountryEquality(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  country.Country
		right country.Country
		want  bool
	}{
		{name: "same named value", left: country.US, right: country.US, want: true},
		{name: "named value and parsed alpha-2", left: country.US, right: parse(t, "us"), want: true},
		{name: "named value and parsed alpha-3", left: country.US, right: parse(t, "USA"), want: true},
		{name: "named value and parsed numeric", left: country.US, right: parse(t, "840"), want: true},
		{name: "different countries", left: country.US, right: country.CA, want: false},
		{name: "country against the zero value", left: country.US, right: country.Country{}, want: false},
		{name: "zero values", left: country.Country{}, right: country.Country{}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.left == tt.right)

			_, found := map[country.Country]bool{tt.left: true}[tt.right]
			assert.Equal(t, tt.want, found, "map key")
		})
	}
}

func TestCountryText(t *testing.T) {
	t.Parallel()

	var (
		_ encoding.TextMarshaler   = country.Country{}
		_ encoding.TextUnmarshaler = &country.Country{}
	)

	tests := []struct {
		name      string
		country   country.Country
		want      string
		readsBack bool
	}{
		{name: "named value", country: country.US, want: "US", readsBack: true},
		{name: "parsed from alpha-3", country: parse(t, "deu"), want: "DE", readsBack: true},
		{name: "parsed from numeric", country: parse(t, "008"), want: "AL", readsBack: true},
		{name: "zero value", country: country.Country{}, want: "", readsBack: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.country.String())

			text, err := tt.country.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(text))

			value, err := tt.country.Value()
			require.NoError(t, err)
			assert.Equal(t, tt.want, value)

			var parsed country.Country

			err = parsed.UnmarshalText(text)
			if !tt.readsBack {
				require.Error(t, err)
				assert.Equal(t, country.Country{}, parsed)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.country, parsed)
		})
	}
}

func TestCountryJSON(t *testing.T) {
	t.Parallel()

	type payload struct {
		Domicile country.Country         `json:"domicile"`
		Weights  map[country.Country]int `json:"weights"`
	}

	in := payload{
		Domicile: country.IE,
		Weights:  map[country.Country]int{country.US: 60, country.JP: 40},
	}

	data, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, `{"domicile":"IE","weights":{"US":60,"JP":40}}`, string(data))

	var out payload

	require.NoError(t, json.Unmarshal(data, &out))
	assert.Equal(t, in, out)

	require.EqualError(t, json.Unmarshal([]byte(`{"domicile":"UK"}`), &out), `invalid country "UK": not an ISO 3166-1 code`)
	assert.Equal(t, in, out, "a failed decode must leave the country unchanged")
}

func TestCountryScan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		src     any
		want    country.Country
		wantErr string
	}{
		{name: "string", src: "DE", want: country.DE},
		{name: "bytes", src: []byte("JP"), want: country.JP},
		{name: "alpha-3 string", src: "che", want: country.CH},
		{name: "numeric bytes", src: []byte("826"), want: country.GB},
		{name: "padded string", src: "FR  ", want: country.FR},
		{name: "null", src: nil, wantErr: "invalid country: NULL"},
		{name: "unsupported type", src: int64(840), wantErr: "invalid country: unsupported type int64"},
		{name: "invalid string", src: "UK", wantErr: `invalid country "UK": not an ISO 3166-1 code`},
		{name: "invalid bytes", src: []byte("EU"), wantErr: `invalid country "EU": not an ISO 3166-1 code`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			before := country.US
			got := before

			err := got.Scan(tt.src)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				assert.Equal(t, before, got, "a failed scan must leave the country unchanged")

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestCountryAllocations does not call t.Parallel: testing.AllocsPerRun
// counts the allocations of every goroutine, so it must not run beside other
// tests.
func TestCountryAllocations(t *testing.T) {
	var (
		weights          = map[country.Country]int{country.US: 1}
		text             = []byte(" usa ")
		stringSource any = "Usa"
		bytesSource  any = []byte("840")
		got          country.Country
		found        bool
	)

	// The runs share got and found, so that the compiler keeps their work.
	// The parsing runs come first and leave got set to US.
	tests := []struct {
		name string
		run  func()
	}{
		{name: "parse alpha-2", run: func() { got, _ = country.NewCountryFromString(" us ") }},
		{name: "parse alpha-3", run: func() { got, _ = country.NewCountryFromString("usa") }},
		{name: "parse numeric", run: func() { got, _ = country.NewCountryFromString("840") }},
		{name: "unmarshal text", run: func() { _ = got.UnmarshalText(text) }},
		{name: "scan string", run: func() { _ = got.Scan(stringSource) }},
		{name: "scan bytes", run: func() { _ = got.Scan(bytesSource) }},
		{name: "copy and compare", run: func() { c := got; found = c == country.US }},
		{name: "hash", run: func() { _, found = weights[got] }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Zero(t, testing.AllocsPerRun(100, tt.run))
		})
	}

	assert.Equal(t, country.US, got)
	assert.True(t, found)
}

func TestCountryNamedValues(t *testing.T) {
	t.Parallel()

	require.Len(t, named, 249, "ISO 3166-1 officially assigns 249 codes")

	for i, c := range named {
		code := c.Alpha2()

		if i > 0 {
			assert.Less(t, named[i-1].Alpha2(), code, "named values must be distinct and in order of alpha-2 code")
		}

		for _, form := range []string{code, c.Alpha3(), fmt.Sprintf("%03d", c.NumericCode())} {
			got, err := country.NewCountryFromString(form)
			if assert.NoError(t, err, "country %s", code) {
				assert.Equal(t, c, got, "country %s read from %q", code, form)
			}
		}

		assert.NotEmpty(t, c.Name(), "country %s", code)
	}
}

// named lists every named value in order of alpha-2 code.
var named = []country.Country{
	country.AD, country.AE, country.AF, country.AG, country.AI, country.AL, country.AM, country.AO, country.AQ, country.AR,
	country.AS, country.AT, country.AU, country.AW, country.AX, country.AZ, country.BA, country.BB, country.BD, country.BE,
	country.BF, country.BG, country.BH, country.BI, country.BJ, country.BL, country.BM, country.BN, country.BO, country.BQ,
	country.BR, country.BS, country.BT, country.BV, country.BW, country.BY, country.BZ, country.CA, country.CC, country.CD,
	country.CF, country.CG, country.CH, country.CI, country.CK, country.CL, country.CM, country.CN, country.CO, country.CR,
	country.CU, country.CV, country.CW, country.CX, country.CY, country.CZ, country.DE, country.DJ, country.DK, country.DM,
	country.DO, country.DZ, country.EC, country.EE, country.EG, country.EH, country.ER, country.ES, country.ET, country.FI,
	country.FJ, country.FK, country.FM, country.FO, country.FR, country.GA, country.GB, country.GD, country.GE, country.GF,
	country.GG, country.GH, country.GI, country.GL, country.GM, country.GN, country.GP, country.GQ, country.GR, country.GS,
	country.GT, country.GU, country.GW, country.GY, country.HK, country.HM, country.HN, country.HR, country.HT, country.HU,
	country.ID, country.IE, country.IL, country.IM, country.IN, country.IO, country.IQ, country.IR, country.IS, country.IT,
	country.JE, country.JM, country.JO, country.JP, country.KE, country.KG, country.KH, country.KI, country.KM, country.KN,
	country.KP, country.KR, country.KW, country.KY, country.KZ, country.LA, country.LB, country.LC, country.LI, country.LK,
	country.LR, country.LS, country.LT, country.LU, country.LV, country.LY, country.MA, country.MC, country.MD, country.ME,
	country.MF, country.MG, country.MH, country.MK, country.ML, country.MM, country.MN, country.MO, country.MP, country.MQ,
	country.MR, country.MS, country.MT, country.MU, country.MV, country.MW, country.MX, country.MY, country.MZ, country.NA,
	country.NC, country.NE, country.NF, country.NG, country.NI, country.NL, country.NO, country.NP, country.NR, country.NU,
	country.NZ, country.OM, country.PA, country.PE, country.PF, country.PG, country.PH, country.PK, country.PL, country.PM,
	country.PN, country.PR, country.PS, country.PT, country.PW, country.PY, country.QA, country.RE, country.RO, country.RS,
	country.RU, country.RW, country.SA, country.SB, country.SC, country.SD, country.SE, country.SG, country.SH, country.SI,
	country.SJ, country.SK, country.SL, country.SM, country.SN, country.SO, country.SR, country.SS, country.ST, country.SV,
	country.SX, country.SY, country.SZ, country.TC, country.TD, country.TF, country.TG, country.TH, country.TJ, country.TK,
	country.TL, country.TM, country.TN, country.TO, country.TR, country.TT, country.TV, country.TW, country.TZ, country.UA,
	country.UG, country.UM, country.US, country.UY, country.UZ, country.VA, country.VC, country.VE, country.VG, country.VI,
	country.VN, country.VU, country.WF, country.WS, country.YE, country.YT, country.ZA, country.ZM, country.ZW,
}
