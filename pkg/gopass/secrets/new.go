package secrets

import (
	"github.com/kpitt/gopass/pkg/gopass"
)

// New creates a new secret.
func New() gopass.Secret { //nolint:ireturn
	return NewKV()
}
