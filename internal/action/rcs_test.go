package action

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/termio"
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
	r1 := gptest.UnsetVars(termio.NameVars...)
	defer r1()
	r2 := gptest.UnsetVars(termio.EmailVars...)
	defer r2()

	un := termio.DetectName(ctx, c)
	ue := termio.DetectEmail(ctx, c)

	name, email := act.getUserData(ctx, "", un, ue)
	assert.Equal(t, "foobar", name)
	assert.Equal(t, "foo.bar@example.org", email)
}
