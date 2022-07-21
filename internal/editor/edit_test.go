package editor

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/stretchr/testify/assert"
)

func TestEdit(t *testing.T) { //nolint:paralleltest
	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithTerminal(ctx, false)

	buf := &bytes.Buffer{}
	out.Stdout = buf
	defer func() {
		out.Stdout = os.Stdout
	}()

	_, err := Invoke(ctx, "true", []byte{})
	assert.Error(t, err)
	buf.Reset()
}
