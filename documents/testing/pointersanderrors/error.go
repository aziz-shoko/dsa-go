package pointersanderrors

import (
	"errors"
	"fmt"
)

// defined errors
var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds ")

type Bitcoin int

type Wallet struct {
	balance Bitcoin
}

// Wallet methods
func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) Withdraw(amount Bitcoin) error {
	if amount > w.balance {
		return ErrInsufficientFunds
	}
	w.balance -= amount
	return nil
}

// Bitcoin methods
func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
