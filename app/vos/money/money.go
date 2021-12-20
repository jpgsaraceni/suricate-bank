// Package money contains functions and methods for creating, adding to and subtracting from
// Money types. Values are expressed in cents (integers) and must always be positive.
package money

import "errors"

type Money struct {
	cents int
}

var (
	ErrNegative         = errors.New("negative values not allowed")
	ErrChangeByZero     = errors.New("cannot add or subtract 0")
	ErrInsuficientFunds = errors.New("subtract amount greater than available amount")
)

func NewMoney(amount int) (Money, error) {
	if amount < 0 {

		return Money{}, ErrNegative
	}

	return Money{cents: amount}, nil
}

func (m Money) Cents() int {
	return m.cents
}

func (m *Money) Add(amount int) error {
	if amount < 0 {

		return ErrNegative
	}

	if amount == 0 {

		return ErrChangeByZero
	}
	m.cents += amount

	return nil
}

func (m *Money) Subtract(amount int) error {
	if amount < 0 {

		return ErrNegative
	}

	if amount == 0 {

		return ErrChangeByZero
	}

	if amount > m.cents {

		return ErrInsuficientFunds
	}
	m.cents -= amount

	return nil
}