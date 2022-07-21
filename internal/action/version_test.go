package action

import (
	"bytes"
	"context"
	"os"
	"testing"

	_ "github.com/kpitt/gopass/internal/backend/crypto"
	_ "github.com/kpitt/gopass/internal/backend/storage"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/tests/gptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestVersion(t *testing.T) {
	t.Parallel()

	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithInteractive(ctx, false)

	act, err := newMock(ctx, u.StoreDir(""))
	require.NoError(t, err)

	buf := &bytes.Buffer{}
	out.Stdout = buf
	stdout = buf
	defer func() {
		out.Stdout = os.Stdout
		stdout = os.Stdout
	}()

	cli.VersionPrinter = func(*cli.Context) {
		out.Printf(ctx, "gopass version 0.0.0-test")
	}

	t.Run("print fixed version", func(t *testing.T) {
		assert.NoError(t, act.Version(gptest.CliCtx(ctx, t)))
	})
}
