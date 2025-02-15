package root

import (
	"context"
	"testing"

	"github.com/fatih/color"
	"github.com/kpitt/gopass/internal/tree"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/tests/gptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Parallel()

	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithHidden(ctx, true)
	color.NoColor = true

	rs, err := createRootStore(ctx, u)
	require.NoError(t, err)

	es, err := rs.List(ctx, tree.INF)
	require.NoError(t, err)
	assert.Equal(t, []string{"foo"}, es)

	sd, err := rs.HasSubDirs(ctx, "foo")
	assert.NoError(t, err)
	assert.False(t, sd)

	str, err := rs.Format(ctx, -1)
	assert.NoError(t, err)
	assert.Equal(t, `gopass
└── foo
`, str)
}
