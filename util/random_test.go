package util

import (
	"testing"
)

func TestRandomInt(t *testing.T) {
	max := int64(100)
	min := int64(0)
	num := RandomInt(max, min)
	if num < min || num > max {
		t.Errorf("RandomInt() = %d; expected value between %d and %d", num, min, max)
	}
}

func TestRandomString(t *testing.T) {
	n := 10
	str := RandomString(n)
	if len(str) != n {
		t.Errorf("RandomString(%d) = %q; expected string of length %d", n, str, n)
	}
}

func TestRandomOwner(t *testing.T) {
	owner := RandomOwner()
	if len(owner) != 6 {
		t.Errorf("RandomOwner() = %q; expected string of length 6", owner)
	}
}

func TestRandomMoney(t *testing.T) {
	money := RandomMoney()
	if money < 0 || money > 10000 {
		t.Errorf("RandomMoney() = %d; expected value between 0 and 10000", money)
	}
}

func TestRandomCurrency(t *testing.T) {
	currency := RandomCurrency()
	if currency != "EUR" && currency != "USD" {
		t.Errorf("RandomCurrency() = %q; expected EUR or USD", currency)
	}
}

func TestRandomEmail(t *testing.T) {
	email := RandomEmail()
	if len(email) != 14+len("com") {
		t.Errorf("RandomEmail() = %q; expected string of length %d", email, 13+len("com"))
	}
	if email[6:7] != "@" {
		t.Errorf("RandomEmail() = %q; expected '@' at index 7", email)
	}
	if email[14:17] != "com" {
		t.Errorf("RandomEmail() = %q; expected 'com' at index 10-12", email)
	}
}
