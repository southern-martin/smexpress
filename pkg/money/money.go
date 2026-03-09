package money

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Amount represents a currency-safe monetary amount.
type Amount struct {
	Value    decimal.Decimal `json:"value"`
	Currency string          `json:"currency"`
}

// New creates a new Amount.
func New(value decimal.Decimal, currency string) Amount {
	return Amount{Value: value, Currency: currency}
}

// FromFloat creates an Amount from a float64.
func FromFloat(value float64, currency string) Amount {
	return Amount{Value: decimal.NewFromFloat(value), Currency: currency}
}

// FromString creates an Amount from a string value.
func FromString(value string, currency string) (Amount, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return Amount{}, fmt.Errorf("parse amount: %w", err)
	}
	return Amount{Value: d, Currency: currency}, nil
}

// Add returns the sum of two amounts (must be same currency).
func (a Amount) Add(b Amount) (Amount, error) {
	if a.Currency != b.Currency {
		return Amount{}, fmt.Errorf("currency mismatch: %s vs %s", a.Currency, b.Currency)
	}
	return Amount{Value: a.Value.Add(b.Value), Currency: a.Currency}, nil
}

// Sub returns the difference of two amounts (must be same currency).
func (a Amount) Sub(b Amount) (Amount, error) {
	if a.Currency != b.Currency {
		return Amount{}, fmt.Errorf("currency mismatch: %s vs %s", a.Currency, b.Currency)
	}
	return Amount{Value: a.Value.Sub(b.Value), Currency: a.Currency}, nil
}

// Mul multiplies the amount by a factor.
func (a Amount) Mul(factor decimal.Decimal) Amount {
	return Amount{Value: a.Value.Mul(factor), Currency: a.Currency}
}

// RoundTo rounds to the given decimal places.
func (a Amount) RoundTo(places int32) Amount {
	return Amount{Value: a.Value.Round(places), Currency: a.Currency}
}

// IsZero checks if the amount is zero.
func (a Amount) IsZero() bool {
	return a.Value.IsZero()
}

// IsPositive checks if the amount is positive.
func (a Amount) IsPositive() bool {
	return a.Value.IsPositive()
}

// String returns a human-readable representation.
func (a Amount) String() string {
	return fmt.Sprintf("%s %s", a.Value.StringFixed(2), a.Currency)
}
