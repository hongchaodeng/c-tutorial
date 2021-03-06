package demo2

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"speter.net/go/exp/math/dec/inf"
)

// International System of Units:
// (Metric prefix)
// https://en.wikipedia.org/wiki/Metric_prefix
//
// Examples:
//   1.5 will be serialized as "1500M" or "1500m"
//   1.5Gi will be serialized as "1536Mi" or "1536mi"
//
//
// ======================================================
//             KUBERNETES CODE BELOW
// ======================================================
//
//
// Quantity is a fixed-point representation of a number.
// It provides convenient marshaling/unmarshaling in JSON and YAML,
// in addition to String() and Int64() accessors.
//
// The serialization format is:
//
// <quantity>        ::= <signedNumber><suffix>
//   (Note that <suffix> may be empty, from the "" case in <decimalSI>.)
// <digit>           ::= 0 | 1 | ... | 9
// <digits>          ::= <digit> | <digit><digits>
// <number>          ::= <digits> | <digits>.<digits> | <digits>. | .<digits>
// <sign>            ::= "+" | "-"
// <signedNumber>    ::= <number> | <sign><number>
// <suffix>          ::= <binarySI> | <decimalExponent> | <decimalSI>
// <binarySI>        ::= Ki | Mi | Gi | Ti | Pi | Ei
//   (International System of units; See: http://physics.nist.gov/cuu/Units/binary.html)
// <decimalSI>       ::= m | "" | k | M | G | T | P | E
//   (Note that 1024 = 1Ki but 1000 = 1k; I didn't choose the capitalization.)
// <decimalExponent> ::= "e" <signedNumber> | "E" <signedNumber>
//
// No matter which of the three exponent forms is used, no quantity may represent
// a number greater than 2^63-1 in magnitude, nor may it have more than 3 decimal
// places. Numbers larger or more precise will be capped or rounded up.
// (E.g.: 0.1m will rounded up to 1m.)
// This may be extended in the future if we require larger or smaller quantities.
//
// When a Quantity is parsed from a string, it will remember the type of suffix
// it had, and will use the same type again when it is serialized.
//
// Before serializing, Quantity will be put in "canonical form".
// This means that Exponent/suffix will be adjusted up or down (with a
// corresponding increase or decrease in Mantissa) such that:
//   a. No precision is lost
//   b. No fractional digits will be emitted
//   c. The exponent (or suffix) is as large as possible.
// The sign will be omitted unless the number is negative.
//
// Examples:
//   1.5 will be serialized as "1500m"
//   1.5Gi will be serialized as "1536Mi"
//
// NOTE: We reserve the right to amend this canonical format, perhaps to
//   allow 1.5 to be canonical.
// TODO: Remove above disclaimer after all bikeshedding about format is over,
//   or after March 2015.
//
// Note that the quantity will NEVER be internally represented by a
// floating point number. That is the whole point of this exercise.
//
// Non-canonical values will still parse as long as they are well formed,
// but will be re-emitted in their canonical form. (So always use canonical
// form, or don't diff.)
//
// This format is intended to make it difficult to use these numbers without
// writing some sort of special handling code in the hopes that that will
// cause implementors to also use a fixed point implementation.
//
// +protobuf=true
// +protobuf.embed=QuantityProto
// +protobuf.options.marshal=false
// +protobuf.options.(gogoproto.goproto_stringer)=false
type Quantity struct {
	// Amount is public, so you can manipulate it if the accessor
	// functions are not sufficient.
	Amount *inf.Dec

	// Change Format at will. See the comment for Canonicalize for
	// more details.
	Format
}

// Format lists the three possible formattings of a quantity.
type Format string

const (
	DecimalExponent = Format("DecimalExponent") // e.g., 12e6
	BinarySI        = Format("BinarySI")        // e.g., 12Mi (12 * 2^20)
	DecimalSI       = Format("DecimalSI")       // e.g., 12M  (12 * 10^6)
)

// MustParse turns the given string into a quantity or panics; for tests
// or others cases where you know the string is valid.
func MustParse(str string) Quantity {
	q, err := ParseQuantity(str)
	if err != nil {
		panic(fmt.Errorf("cannot parse '%v': %v", str, err))
	}
	return *q
}

// Scale is used for getting and setting the base-10 scaled value.
// Base-2 scales are omitted for mathematical simplicity.
// See Quantity.ScaledValue for more details.
type Scale int

const (
	Nano  Scale = -9
	Micro Scale = -6
	Milli Scale = -3
	Kilo  Scale = 3
	Mega  Scale = 6
	Giga  Scale = 9
	Tera  Scale = 12
	Peta  Scale = 15
	Exa   Scale = 18
)

const (
	// splitREString is used to separate a number from its suffix; as such,
	// this is overly permissive, but that's OK-- it will be checked later.
	splitREString = "^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$"
)

var (
	// splitRE is used to get the various parts of a number.
	splitRE = regexp.MustCompile(splitREString)

	// Errors that could happen while parsing a string.
	ErrFormatWrong = errors.New("quantities must match the regular expression '" + splitREString + "'")
	ErrNumeric     = errors.New("unable to parse numeric part of quantity")
	ErrSuffix      = errors.New("unable to parse quantity's suffix")

	// Commonly needed big.Int values-- treat as read only!
	bigTen      = big.NewInt(10)
	bigZero     = big.NewInt(0)
	bigOne      = big.NewInt(1)
	bigThousand = big.NewInt(1000)
	big1024     = big.NewInt(1024)

	// Commonly needed inf.Dec values-- treat as read only!
	decZero      = inf.NewDec(0, 0)
	decOne       = inf.NewDec(1, 0)
	decMinusOne  = inf.NewDec(-1, 0)
	decThousand  = inf.NewDec(1000, 0)
	dec1024      = inf.NewDec(1024, 0)
	decMinus1024 = inf.NewDec(-1024, 0)

	// Largest (in magnitude) number allowed.
	maxAllowed = inf.NewDec((1<<63)-1, 0) // == max int64

	// The maximum value we can represent milli-units for.
	// Compare with the return value of Quantity.Value() to
	// see if it's safe to use Quantity.MilliValue().
	MaxMilliValue = int64(((1 << 63) - 1) / 1000)
)

// ParseQuantity turns str into a Quantity, or returns an error.
func ParseQuantity(str string) (*Quantity, error) {
	parts := splitRE.FindStringSubmatch(strings.TrimSpace(str))
	// regexp returns are entire match, followed by an entry for each () section.
	if len(parts) != 3 {
		return nil, ErrFormatWrong
	}

	amount := new(inf.Dec)
	if _, ok := amount.SetString(parts[1]); !ok {
		return nil, ErrNumeric
	}

	base, exponent, format, ok := quantitySuffixer.interpret(suffix(parts[2]))
	if !ok {
		return nil, ErrSuffix
	}

	// So that no one but us has to think about suffixes, remove it.
	if base == 10 {
		amount.SetScale(amount.Scale() + Scale(exponent).infScale())
	} else if base == 2 {
		// numericSuffix = 2 ** exponent
		numericSuffix := big.NewInt(1).Lsh(bigOne, uint(exponent))
		ub := amount.UnscaledBig()
		amount.SetUnscaledBig(ub.Mul(ub, numericSuffix))
	}

	// Cap at min/max bounds.
	sign := amount.Sign()
	if sign == -1 {
		amount.Neg(amount)
	}

	// This rounds non-zero values up to the minimum representable value, under the theory that
	// if you want some resources, you should get some resources, even if you asked for way too small
	// of an amount.  Arguably, this should be inf.RoundHalfUp (normal rounding), but that would have
	// the side effect of rounding values < .5n to zero.
	if v, ok := amount.Unscaled(); v != int64(0) || !ok {
		amount.Round(amount, Nano.infScale(), inf.RoundUp)
	}

	// The max is just a simple cap.
	if amount.Cmp(maxAllowed) > 0 {
		amount.Set(maxAllowed)
	}
	if format == BinarySI && amount.Cmp(decOne) < 0 && amount.Cmp(decZero) > 0 {
		// This avoids rounding and hopefully confusion, too.
		format = DecimalSI
	}
	if sign == -1 {
		amount.Neg(amount)
	}

	return &Quantity{amount, format}, nil
}

// Value returns the value of q; any fractional part will be lost.
func (q *Quantity) Value() int64 {
	if q.Amount == nil {
		return 0
	}
	tmp := &inf.Dec{}
	return tmp.Round(q.Amount, 0, inf.RoundUp).UnscaledBig().Int64()
	// return q.ScaledValue(0)
}

// MilliValue returns the value of ceil(q * 1000); this could overflow an int64;
// if that's a concern, call Value() first to verify the number is small enough.
func (q *Quantity) MilliValue() int64 {
	if q.Amount == nil {
		return 0
	}
	// Numbers larger or more precise will be capped or rounded up.
	// (E.g.: 0.1m will rounded up to 1m.)
	tmp := &inf.Dec{}
	return tmp.Round(tmp.Mul(q.Amount, decThousand), 0, inf.RoundUp).UnscaledBig().Int64()
	// return scaledValue(q.Amount.UnscaledBig(), int(q.Amount.Scale()), 3)
}

// infScale adapts a Scale value to an inf.Scale value.
func (s Scale) infScale() inf.Scale {
	return inf.Scale(-s) // inf.Scale is upside-down
}
