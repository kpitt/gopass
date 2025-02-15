package action

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/gopass/secrets"
	"github.com/kpitt/gopass/tests/gptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestFind(t *testing.T) { //nolint:paralleltest
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithTerminal(ctx, false)

	act, err := newMock(ctx, u.StoreDir(""))
	require.NoError(t, err)
	require.NotNil(t, act)

	buf := &bytes.Buffer{}
	out.Stdout = buf
	stdout = buf
	defer func() {
		stdout = os.Stdout
		out.Stdout = os.Stdout
	}()
	color.NoColor = true

	actName := "action.test"
	if runtime.GOOS == "windows" {
		actName = "action.test.exe"
	}

	// find
	c := gptest.CliCtx(ctx, t)
	if err := act.Find(c); err == nil || err.Error() != fmt.Sprintf("Usage: %s find <pattern>", actName) {
		t.Errorf("Should fail: %s", err)
	}

	// find fo (with fuzzy search)
	c = gptest.CliCtxWithFlags(ctx, t, nil, "fo")
	assert.NoError(t, act.FindFuzzy(c, "fo"))
	assert.Contains(t, strings.TrimSpace(buf.String()), "Found exact match in \"foo\"\nsecret")
	buf.Reset()

	// find fo (no fuzzy search)
	c = gptest.CliCtxWithFlags(ctx, t, nil, "fo")
	assert.NoError(t, act.Find(c))
	assert.Equal(t, strings.TrimSpace(buf.String()), "foo")
	buf.Reset()

	// find yo
	c = gptest.CliCtx(ctx, t)
	assert.Error(t, act.FindFuzzy(c, "yo"))
	buf.Reset()

	// add some secrets
	sec := &secrets.Plain{}
	sec.SetPassword("foo")
	sec.WriteString("bar")
	assert.NoError(t, act.Store.Set(ctx, "bar/baz", sec))
	assert.NoError(t, act.Store.Set(ctx, "bar/zab", sec))
	buf.Reset()

	// find bar
	c = gptest.CliCtx(ctx, t)
	assert.NoError(t, act.FindFuzzy(c, "bar"))
	assert.Equal(t, "bar/baz\nbar/zab", strings.TrimSpace(buf.String()))
	buf.Reset()

	// find w/o callback
	c = gptest.CliCtx(ctx, t)
	assert.NoError(t, act.find(ctx, c, "foo", nil, false))
	assert.Equal(t, "foo", strings.TrimSpace(buf.String()))
	buf.Reset()

	// findSelection w/o callback
	c = gptest.CliCtx(ctx, t)
	assert.Error(t, act.findSelection(ctx, c, []string{"foo", "bar"}, "fo", nil))

	// findSelection w/o options
	c = gptest.CliCtx(ctx, t)
	assert.Error(t, act.findSelection(ctx, c, nil, "fo", func(_ context.Context, _ *cli.Context, _ string, _ bool) error { return nil }))
}
