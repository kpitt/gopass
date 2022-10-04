package cui

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/kpitt/gopass/internal/backend/crypto/plain"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAskForPrivateKey(t *testing.T) { //nolint:paralleltest
	buf := &bytes.Buffer{}
	Stdout = buf
	defer func() {
		Stdout = os.Stdout
	}()

	ctx := context.Background()

	ctx = ctxutil.WithAlwaysYes(ctx, true)
	key, err := AskForPrivateKey(ctx, plain.New(), "test")
	require.NoError(t, err)
	assert.Equal(t, "0xDEADBEEF", key)
	buf.Reset()
}

type fakeMountPointer []string

func (f fakeMountPointer) MountPoints() []string {
	return f
}

func TestAskForStore(t *testing.T) { //nolint:paralleltest
	ctx := context.Background()

	// test non-interactive
	ctx = ctxutil.WithInteractive(ctx, false)
	assert.Equal(t, "", AskForStore(ctx, fakeMountPointer{"foo", "bar"}))

	// test interactive
	ctx = ctxutil.WithInteractive(ctx, true)
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	assert.Equal(t, "", AskForStore(ctx, fakeMountPointer{"foo", "bar"}))

	// test zero mps
	assert.Equal(t, "", AskForStore(ctx, fakeMountPointer{}))

	// test one mp
	assert.Equal(t, "", AskForStore(ctx, fakeMountPointer{"foo"}))

	// test two mps
	assert.Equal(t, "", AskForStore(ctx, fakeMountPointer{"foo", "bar"}))
}

func TestSorted(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []string{"a", "b", "c"}, sorted([]string{"c", "a", "b"}))
}
