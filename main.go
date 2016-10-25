package main

import (
	"fmt"
)

type Bigint struct {
	digits []int
	minus bool
}

func New(n ...int) *Bigint {
	bigint := &Bigint{}

	if len(n) > 0 {
		bigint.digits = make([]int, n[0])
	}

	return bigint
}

func (b Bigint) IsMinus() bool {
	return b.minus == true
}

func (b *Bigint) Minus() *Bigint {
	b.minus = true

	return b
}

func (b Bigint) Digits() []int {
	return b.digits
} 

func (b *Bigint) Set(digits []int) *Bigint {
	b.digits = digits

	return b
}

func (b Bigint) Len() int {
	return len(b.digits)
}

func (b Bigint) String() string {
	result := ""
	hadDigit := false

	if b.IsMinus() {
		result = "-"
	}

	for _, d := range b.Digits() {
		if d != 0 {
			hadDigit = true
		}

		if !hadDigit && d == 0 {
			continue
		}

		result += fmt.Sprintf("%d", d)
	}

	return result
}

func (b *Bigint) Fill(n int) *Bigint {
	d := n - b.Len()

	if d > 0 {
		digits := make([]int, n)

		for i := d; i < n; i++ {
			digits[i] = b.digits[i - d]
		}

		b.digits = digits
	}

	return b
}

func (b *Bigint) IsEqual(a *Bigint) bool {
	aLen := a.Len()
	bLen := b.Len()

	if a.IsMinus() != b.IsMinus() {
		return false
	}

	if aLen != bLen {
		return false
	}

	for i := 0; i < aLen; i++ {
		if a.Digits()[i] != b.Digits()[i] {
			return false
		}
	}

	return true
}

func (b *Bigint) IsLesser(a *Bigint) bool {
	if b.IsEqual(a) {
		return false
	}

	aLen := a.Len()
	bLen := b.Len()

	if a.IsMinus() != b.IsMinus() {
		if b.IsMinus() {
			return true
		}

		return false
	}

	if aLen != bLen {
		if aLen > bLen {
			return !a.IsMinus()
		}

		return !b.IsMinus()
	}

	for i := 0; i < aLen; i++ {
		if a.IsMinus() {
			if a.Digits()[i] > b.Digits()[i] {
				return false
			}
		} else {
			if a.Digits()[i] < b.Digits()[i] {
				return false
			}
		}
	}

	return true
}

func (b *Bigint) IsGreater(a *Bigint) bool {
	if b.IsEqual(a) {
		return false
	}

	aLen := a.Len()
	bLen := b.Len()

	if a.IsMinus() != b.IsMinus() {
		if a.IsMinus() {
			return true
		}

		return false
	}

	if aLen != bLen {
		if aLen < bLen {
			return !b.IsMinus()
		}

		return !a.IsMinus()
	}

	for i := 0; i < aLen; i++ {
		if a.IsMinus() {
			if a.Digits()[i] < b.Digits()[i] {
				return false
			}
		} else {
			if a.Digits()[i] > b.Digits()[i] {
				return false
			}
		}
	}

	return true
}

func (b *Bigint) Sub(a *Bigint) *Bigint {
	c := New()

	switch {
	case a.IsMinus() && b.IsMinus():
		a.minus = false
		c = b.Add(a)
		a.minus = true
		
		if a.IsGreater(b) {
			c.minus = true
		}

		return c
	case a.IsMinus() && !b.IsMinus():
		a.minus = false
		c = b.Add(a)
		a.minus = true

		return c
	case !a.IsMinus() && b.IsMinus():
		b.minus = false
		c = b.Add(a)
		b.minus = true
		c.minus = true

		return c
	}

	aLen := a.Len()
	bLen := b.Len()

	maxLen := aLen

	if bLen > aLen {
		maxLen = bLen
	}

	a.Fill(maxLen)
	b.Fill(maxLen)

	aDigits := a.Digits()
	bDigits := b.Digits()
	cDigits := make([]int, maxLen)

	if b.IsLesser(a) {
		aDigits, bDigits = bDigits, aDigits
		c.minus = true
	}

	inMemory := 0

	for i := maxLen - 1; i >= 0; i-- {
		d := bDigits[i] - aDigits[i]

		if d < 0 {
			d = d + 10

			if inMemory < 0 {
				d--
				inMemory++
			}

			inMemory--
		} else {
			if inMemory < 0 {
				d--
				inMemory++
			}

			if d < 0 {
				d = d + 10
				inMemory--
			}
		}

		cDigits[i] = d
	}

	if inMemory < 0 {
		cDigits[0]--
	}

	c.digits = cDigits

	return c
}

func (b *Bigint) Add(a *Bigint) *Bigint {
	c := New()

	switch {
	case a.IsMinus() && b.IsMinus():
		c.minus = true
	case a.IsMinus() && !b.IsMinus():
		a.minus = false
		c = b.Sub(a)
		a.minus = true

		return c
	case !a.IsMinus() && b.IsMinus():
		b.minus = false
		c = b.Sub(a)
		b.minus = true

		if c.IsMinus() {
			c.minus = false
		}

		return c
	}

	aLen := a.Len()
	bLen := b.Len()

	maxLen := aLen

	if bLen > aLen {
		maxLen = bLen
	}

	digits := make([]int, maxLen + 1)

	a.Fill(maxLen)
	b.Fill(maxLen)

	inMemory := 0

	for i := maxLen - 1; i >= 0; i-- {
		d := a.Digits()[i] + b.Digits()[i]

		if d > 9 {
			d = d - 10

			if inMemory > 0 {
				d++
				inMemory--
			}

			inMemory++
		} else {
			if inMemory > 0 {
				d++
				inMemory--
			}

			if d > 9 {
				d = d - 10
				inMemory++
			}
		}

		digits[i + 1] = d
	}

	if inMemory > 0 {
		digits[0] = inMemory
	} else {
		digits = digits[1:]
	}

	c.Set(digits)

	return c
}

func main() {
	a := New()
	a.Set([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}).Minus()

	b := New()
	b.Set([]int{1, 2, 3, 4, 5}).Minus()

	c := b.Sub(a)

	fmt.Println(c)
}