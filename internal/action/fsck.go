package action

import (
	"github.com/kpitt/gopass/internal/action/exit"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/internal/store/leaf"
	"github.com/kpitt/gopass/internal/tree"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/termio"
	"github.com/urfave/cli/v2"
)

// Fsck checks the store integrity.
func (s *Action) Fsck(c *cli.Context) error {
	_ = s.rem.Reset("fsck")

	ctx := ctxutil.WithGlobalFlags(c)
	if c.IsSet("decrypt") {
		ctx = leaf.WithFsckDecrypt(ctx, c.Bool("decrypt"))
	}

	out.Printf(ctx, "Checking password store integrity...\n")
	// make sure config is in the right place.
	// we may have loaded it from one of the fallback locations.
	if err := s.cfg.Save(); err != nil {
		return exit.Error(exit.Config, err, "failed to save config: %s", err)
	}

	// display progress bar.
	t, err := s.Store.Tree(ctx)
	if err != nil {
		return exit.Error(exit.Unknown, err, "failed to list stores: %s", err)
	}

	pwList := t.List(tree.INF)

	mounts := s.Store.MountPoints()
	steps := len(pwList)*2 + len(mounts) + 1
	bar := termio.NewProgressBar("Checking storage backend", int64(steps))
	bar.Hidden = ctxutil.IsHidden(ctx)
	ctx = ctxutil.WithProgressCallback(ctx, func(msg string) {
		bar.SetText(msg)
		bar.Inc()
	})
	ctx = out.AddPrefix(ctx, "\n")

	// the main work in done by the sub stores.
	if err := s.Store.Fsck(ctx, c.Args().Get(0)); err != nil {
		return exit.Error(exit.Fsck, err, "fsck found errors: %s", err)
	}
	bar.Done()

	return nil
}
