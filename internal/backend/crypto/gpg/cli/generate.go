package cli

import (
	"context"

	"github.com/kpitt/gopass/internal/backend"
)

// GenerateIdentity will create a new GPG keypair in batch mode.
func (g *GPG) GenerateIdentity(ctx context.Context, passphrase string) error {
	return backend.ErrNotSupported
}
