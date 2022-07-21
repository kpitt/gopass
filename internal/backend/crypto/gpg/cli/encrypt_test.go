package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/kpitt/gopass/internal/backend"
	"github.com/kpitt/gopass/internal/backend/crypto/gpg/gpgconf"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/fsutil"
	"github.com/kpitt/gopass/tests/can"
	"github.com/kpitt/gopass/tests/gptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) { //nolint:paralleltest
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	u := gptest.NewGUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithTerminal(ctx, false)
	ctx = backend.WithCryptoBackend(ctx, backend.GPGCLI)

	g, err := New(ctx, Config{
		Umask: fsutil.Umask(),
		Args:  gpgconf.GPGOpts(),
	})
	require.NoError(t, err)

	// import keys so GPG4Win can find them
	el := can.EmbeddedKeyRing()
	for _, e := range el {
		buf := &bytes.Buffer{}
		require.NoError(t, e.Serialize(buf))

		require.NoError(t, g.ImportPublicKey(ctx, buf.Bytes()))
	}

	plaintext := []byte("plaintext")
	ciphertext, err := g.Encrypt(ctx, plaintext, []string{can.KeyID()})
	require.NoError(t, err)

	plaintext2, err := g.Decrypt(ctx, ciphertext)
	require.NoError(t, err)

	assert.Equal(t, plaintext, plaintext2)
}
