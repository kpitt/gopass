package action

import (
	"github.com/kpitt/gopass/internal/action/exit"
	"github.com/kpitt/gopass/internal/audit"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/internal/tree"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/debug"
	"github.com/urfave/cli/v2"
)

// Audit validates passwords against common flaws.
func (s *Action) Audit(c *cli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	expiry := c.Int("expiry")
	if expiry > 0 {
		out.Print(ctx, "Auditing password expiration ...")
	} else {
		_ = s.rem.Reset("audit")
		out.Print(ctx, "Auditing passwords for common flaws ...")
	}

	t, err := s.Store.Tree(ctx)
	if err != nil {
		return exit.Error(exit.List, err, "Failed to get store tree: %s", err)
	}

	if filter := c.Args().First(); filter != "" {
		subtree, err := t.FindFolder(filter)
		if err != nil {
			return exit.Error(exit.Unknown, err, "Failed to find subtree: %s", err)
		}
		debug.Log("subtree for %q: %+v", filter, subtree)
		t = subtree
	}
	list := t.List(tree.INF)

	if len(list) < 1 {
		out.Printf(ctx, "No secrets found")

		return nil
	}

	return audit.Batch(ctx, list, s.Store, expiry)
}
