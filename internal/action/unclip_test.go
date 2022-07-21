package action

import (
	"bytes"
	"context"
	"os"
	"testing"

	_ "github.com/kpitt/gopass/internal/backend/crypto"
	_ "github.com/kpitt/gopass/internal/backend/storage"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/tests/gptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnclip(t *testing.T) { //nolint:paralleltest
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	buf := &bytes.Buffer{}
	out.Stdout = buf
	stdout = buf
	defer func() {
		out.Stdout = os.Stdout
		stdout = os.Stdout
	}()

	ctx := context.Background()
	act, err := newMock(ctx, u.StoreDir(""))
	require.NoError(t, err)
	require.NotNil(t, act)

	t.Run("unlcip should fail", func(t *testing.T) {
		assert.Error(t, act.Unclip(gptest.CliCtxWithFlags(ctx, t, map[string]string{"timeout": "0"})))
	})
}
