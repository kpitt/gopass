package leaf

import (
	"context"

	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/internal/store"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/debug"
	"github.com/kpitt/gopass/pkg/gopass"
	"github.com/kpitt/gopass/pkg/gopass/secrets"
	"github.com/kpitt/gopass/pkg/gopass/secrets/secparse"
)

// Get returns the plaintext of a single key.
func (s *Store) Get(ctx context.Context, name string) (gopass.Secret, error) {
	p := s.Passfile(name)

	ciphertext, err := s.storage.Get(ctx, p)
	if err != nil {
		debug.Log("File %s not found: %s", p, err)

		return nil, store.ErrNotFound
	}

	content, err := s.crypto.Decrypt(ctx, ciphertext)
	if err != nil {
		out.Errorf(ctx, "Decryption failed: %s\n%s", err, string(content))

		return nil, store.ErrDecrypt
	}

	if !ctxutil.IsShowParsing(ctx) {
		return secrets.ParsePlain(content), nil
	}

	return secparse.Parse(content)
}
