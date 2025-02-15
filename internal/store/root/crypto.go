package root

import (
	"context"

	"github.com/kpitt/gopass/internal/backend"
	"github.com/kpitt/gopass/pkg/debug"
)

// Crypto returns the crypto backend.
func (r *Store) Crypto(ctx context.Context, name string) backend.Crypto {
	sub, _ := r.getStore(name)
	if !sub.Valid() {
		debug.Log("Sub-Store not found for %s. Returning nil crypto backend", name)

		return nil
	}

	return sub.Crypto()
}
