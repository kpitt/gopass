package action

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/tests/gptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGit(t *testing.T) { //nolint:paralleltest
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithInteractive(ctx, false)

	act, err := newMock(ctx, u.StoreDir(""))
	require.NoError(t, err)
	require.NotNil(t, act)

	buf := &bytes.Buffer{}
	out.Stdout = buf
	stdout = buf
	defer func() {
		out.Stdout = os.Stdout
		stdout = os.Stdout
	}()

	// git init
	c := gptest.CliCtxWithFlags(ctx, t, map[string]string{"name": "foobar", "email": "foo.bar@example.org"})
	assert.NoError(t, act.RCSInit(c))
	buf.Reset()

	// getUserData
	name, email := act.getUserData(ctx, "", "", "")
	assert.Equal(t, "0xDEADBEEF", name)
	assert.Equal(t, "0xDEADBEEF", email)
}
