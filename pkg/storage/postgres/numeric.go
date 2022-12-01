package postgres

import (
	"database/sql"
	"database/sql/driver"
	"math/big"

	"github.com/pkg/errors"
)

// Numeric -
type Numeric big.Int

// NewNumeric -
func NewNumeric(x *big.Int) *Numeric {
	return (*Numeric)(x)
}

func FromInt64(x int64) *Numeric {
	return new(Numeric).FromBigInt(big.NewInt(x))
}

func FromString(x string) (*Numeric, error) {
	a := big.NewInt(0)
	b, ok := a.SetString(x, 10)

	if !ok {
		return nil, errors.Errorf("cannot create Bigint from string")
	}

	return NewNumeric(b), nil
}

func (b *Numeric) Value() (driver.Value, error) {
	return (*big.Int)(b).String(), nil
}

func (b *Numeric) Scan(value interface{}) error {
	var i sql.NullString

	if err := i.Scan(value); err != nil {
		return err
	}

	if _, ok := (*big.Int)(b).SetString(i.String, 10); ok {
		return nil
	}

	return errors.Errorf("Error converting type %T into Bigint", value)
}

func (b *Numeric) toBigInt() *big.Int {
	return (*big.Int)(b)
}

func (b *Numeric) Sub(x *Numeric) *Numeric {
	return (*Numeric)(big.NewInt(0).Sub(b.toBigInt(), x.toBigInt()))
}

func (b *Numeric) Add(x *Numeric) *Numeric {
	return (*Numeric)(big.NewInt(0).Add(b.toBigInt(), x.toBigInt()))
}

func (b *Numeric) Mul(x *Numeric) *Numeric {
	return (*Numeric)(big.NewInt(0).Mul(b.toBigInt(), x.toBigInt()))
}

func (b *Numeric) Div(x *Numeric) *Numeric {
	return (*Numeric)(big.NewInt(0).Div(b.toBigInt(), x.toBigInt()))
}

func (b *Numeric) Neg() *Numeric {
	return (*Numeric)(big.NewInt(0).Neg(b.toBigInt()))
}

func (b *Numeric) ToUInt64() uint64 {
	return b.toBigInt().Uint64()
}

func (b *Numeric) ToInt64() int64 {
	return b.toBigInt().Int64()
}

// same as NewBigint()
func (b *Numeric) FromBigInt(x *big.Int) *Numeric {
	return (*Numeric)(x)
}

func (b *Numeric) String() string {
	return b.toBigInt().String()
}
