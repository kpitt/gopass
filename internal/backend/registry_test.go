package backend_test

import (
	"context"
	"testing"

	"github.com/kpitt/gopass/internal/backend"
	_ "github.com/kpitt/gopass/internal/backend/crypto"
	"github.com/kpitt/gopass/internal/backend/crypto/plain"
	_ "github.com/kpitt/gopass/internal/backend/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeCryptoLoader struct{}

func (l fakeCryptoLoader) New(context.Context) (backend.Crypto, error) {
	return plain.New(), nil
}

func (l fakeCryptoLoader) String() string {
	return "fakecryptoloader"
}

func (l fakeCryptoLoader) Handles(_ context.Context, _ backend.Storage) error {
	return nil
}

func (l fakeCryptoLoader) Priority() int {
	return 1
}

func TestCryptoLoader(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	backend.CryptoRegistry.Register(backend.Plain, "plain", fakeCryptoLoader{})
	c, err := backend.NewCrypto(ctx, backend.Plain)
	require.NoError(t, err)
	assert.Equal(t, c.Name(), "plain")
}
