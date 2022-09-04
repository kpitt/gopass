package xkcdgen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	t.Parallel()

	pw := Random()
	if len(pw) < 4 {
		t.Errorf("too short")
	}

	if len(strings.Fields(pw)) < 4 {
		t.Errorf("too few words")
	}
}

func TestRandomLengthDelim(t *testing.T) {
	t.Parallel()

	sep := "-"
	pw := RandomLengthDelim(10, sep)
	words := strings.Split(pw, sep)
	assert.Len(t, words, 10, "wrong number of words")
}
