package pwgen

import (
	"bytes"
	"fmt"

	"github.com/kpitt/gopass/pkg/debug"
	"github.com/muesli/crunchy"
)

// ErrCrypticInvalid is returned when a password is invalid.
var ErrCrypticInvalid = fmt.Errorf("password does not satisfy all validators")

// Cryptic is a generator for hard-to-remember passwords as required by (too)
// many sites. Prefer memorable or xkcd-style passwords, if possible.
type Cryptic struct {
	Chars      string
	Length     int
	MaxTries   int
	Validators []func(string) error
}

// NewCryptic creates a new generator with sane defaults.
func NewCryptic(length int, symbols bool) *Cryptic {
	if length < 1 {
		length = 16
	}

	chars := Digits + Upper + Lower

	if symbols {
		chars += Syms
	}

	return &Cryptic{
		Chars:    chars,
		Length:   length,
		MaxTries: 64,
	}
}

// NewCrypticWithAllClasses returns a password generator that generates passwords
// containing all available character classes.
func NewCrypticWithAllClasses(length int, symbols bool) *Cryptic {
	c := NewCryptic(length, symbols)
	c.Validators = append(c.Validators, func(pw string) error {
		if containsAllClasses(pw, c.Chars) {
			return nil
		}

		return fmt.Errorf("password does not contain all classes: %w", ErrCrypticInvalid)
	})

	return c
}

// NewCrypticWithCrunchy returns a password generators that only returns a
// password if it's successfully validated with crunchy.
func NewCrypticWithCrunchy(length int, symbols bool) *Cryptic {
	c := NewCryptic(length, symbols)
	c.MaxTries = 3
	validator := crunchy.NewValidator()
	c.Validators = append(c.Validators, validator.Check)

	return c
}

// Password returns a single password from the generator.
func (c *Cryptic) Password() string {
	round := 0
	maxFn := func() bool {
		round++

		if c.MaxTries < 1 {
			return false
		}

		if c.MaxTries == 0 && round >= 64 {
			return true
		}

		if round > c.MaxTries {
			return true
		}

		return false
	}

	for {
		if maxFn() {
			debug.Log("failed to generate password after %d rounds", round)

			return ""
		}

		if pw := c.randomString(); c.isValid(pw) {
			return pw
		}
	}
}

func (c *Cryptic) isValid(pw string) bool {
	for _, v := range c.Validators {
		if err := v(pw); err != nil {
			debug.Log("failed to validate: %s", err)

			return false
		}
	}

	return true
}

func (c *Cryptic) randomString() string {
	pw := &bytes.Buffer{}
	for pw.Len() < c.Length {
		_ = pw.WriteByte(c.Chars[randomInteger(len(c.Chars))])
	}

	return pw.String()
}
