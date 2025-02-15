package action

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/kpitt/gopass/internal/backend"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/tests/can"
	"github.com/kpitt/gopass/tests/gptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFsck(t *testing.T) { //nolint:paralleltest
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithTerminal(ctx, false)
	act, err := newMock(ctx, u.StoreDir(""))
	require.NoError(t, err)
	require.NotNil(t, act)

	buf := &bytes.Buffer{}
	out.Stdout = buf
	out.Stderr = buf
	stdout = buf
	defer func() {
		stdout = os.Stdout
		out.Stdout = os.Stdout
		out.Stderr = os.Stderr
	}()
	color.NoColor = true

	// fsck
	assert.NoError(t, act.Fsck(gptest.CliCtx(ctx, t)))
	output := strings.TrimSpace(buf.String())
	assert.Contains(t, output, "Checking password store integrity...")
	assert.Contains(t, output, "Extra recipients on foo: [0xFEEDBEEF]")
	buf.Reset()

	// fsck (hidden)
	assert.NoError(t, act.Fsck(gptest.CliCtx(ctxutil.WithHidden(ctx, true), t)))
	output = strings.TrimSpace(buf.String())
	assert.NotContains(t, output, "Checking password store integrity...")
	assert.NotContains(t, output, "Extra recipients on foo: [0xFEEDBEEF]")
	buf.Reset()

	// fsck --decrypt
	assert.NoError(t, act.Fsck(gptest.CliCtxWithFlags(ctx, t, map[string]string{"decrypt": "true"})))
	output = strings.TrimSpace(buf.String())
	assert.Contains(t, output, "Checking password store integrity...")
	assert.Contains(t, output, "Extra recipients on foo: [0xFEEDBEEF]")
	buf.Reset()

	// fsck fo
	assert.NoError(t, act.Fsck(gptest.CliCtx(ctx, t, "fo")))
	output = strings.TrimSpace(buf.String())
	assert.Contains(t, output, "Checking password store integrity...")
	assert.Contains(t, output, "Extra recipients on foo: [0xFEEDBEEF]")
	buf.Reset()
}

func TestFsckGpg(t *testing.T) { //nolint:paralleltest
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	u := gptest.NewGUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithTerminal(ctx, false)
	ctx = backend.WithCryptoBackend(ctx, backend.GPGCLI)

	act, err := newMock(ctx, u.StoreDir(""))
	require.NoError(t, err)
	require.NotNil(t, act)

	buf := &bytes.Buffer{}
	out.Stdout = buf
	out.Stderr = buf
	stdout = buf
	defer func() {
		stdout = os.Stdout
		out.Stdout = os.Stdout
		out.Stderr = os.Stderr
	}()
	color.NoColor = true

	sub, err := act.Store.GetSubStore("")
	require.NoError(t, err)
	require.NoError(t, sub.ImportMissingPublicKeys(ctx, can.KeyID()))

	// generate foo
	c := gptest.CliCtx(ctx, t, "foo", "24")
	assert.NoError(t, act.Generate(c))
	buf.Reset()

	// fsck
	assert.NoError(t, act.Fsck(gptest.CliCtx(ctx, t)))
	output := strings.TrimSpace(buf.String())
	assert.Contains(t, output, "Checking password store integrity...")
	buf.Reset()

	// fsck (hidden)
	assert.NoError(t, act.Fsck(gptest.CliCtx(ctxutil.WithHidden(ctx, true), t)))
	output = strings.TrimSpace(buf.String())
	assert.NotContains(t, output, "Checking password store integrity...")
	assert.NotContains(t, output, "Extra recipients on foo: [0xFEEDBEEF]")
	buf.Reset()

	// fsck --decrypt
	assert.NoError(t, act.Fsck(gptest.CliCtxWithFlags(ctx, t, map[string]string{"decrypt": "true"})))
	output = strings.TrimSpace(buf.String())
	assert.Contains(t, output, "Checking password store integrity...")
	buf.Reset()

	// fsck fo
	assert.NoError(t, act.Fsck(gptest.CliCtx(ctx, t, "fo")))
	output = strings.TrimSpace(buf.String())
	assert.Contains(t, output, "Checking password store integrity...")
	buf.Reset()
}
