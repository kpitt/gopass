package leaf

import (
	"context"
	"fmt"

	"github.com/kpitt/gopass/internal/backend"
	"github.com/kpitt/gopass/internal/store"
	"github.com/kpitt/gopass/pkg/debug"
	"github.com/kpitt/gopass/pkg/gopass"
	"github.com/kpitt/gopass/pkg/gopass/secrets/secparse"
)

// GitInit initializes the git storage.
func (s *Store) GitInit(ctx context.Context) error {
	// The desired storage type for `GitInit` is always GitFS.
	storage, err := backend.InitStorage(ctx, backend.GitFS, s.path)
	if err != nil {
		return err
	}
	s.storage = storage

	return nil
}

// ListRevisions will list all revisions for a secret.
func (s *Store) ListRevisions(ctx context.Context, name string) ([]backend.Revision, error) {
	p := s.Passfile(name)

	return s.storage.Revisions(ctx, p)
}

// GetRevision will retrieve a single revision from the backend.
func (s *Store) GetRevision(ctx context.Context, name, revision string) (gopass.Secret, error) {
	p := s.Passfile(name)
	ciphertext, err := s.storage.GetRevision(ctx, p, revision)
	if err != nil {
		return nil, fmt.Errorf("failed to get ciphertext of %q@%q: %w", name, revision, err)
	}

	content, err := s.crypto.Decrypt(ctx, ciphertext)
	if err != nil {
		debug.Log("Decryption failed: %s", err)

		return nil, store.ErrDecrypt
	}

	sec, err := secparse.Parse(content)
	if err != nil {
		debug.Log("Failed to parse YAML: %s", err)
	}

	return sec, nil
}
