package action

import (
	"fmt"

	"github.com/kpitt/gopass/internal/action/exit"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/termio"
	"github.com/urfave/cli/v2"
)

// Move the content from one secret to another.
func (s *Action) Move(c *cli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

	if c.Args().Len() != 2 {
		return exit.Error(exit.Usage, nil, "Usage: %s mv old-path new-path", s.Name)
	}

	from := c.Args().Get(0)
	to := c.Args().Get(1)

	if !c.Bool("force") {
		if s.Store.Exists(ctx, to) && !termio.AskForConfirmation(ctx, fmt.Sprintf("%s already exists. Overwrite it?", to)) {
			return exit.Error(exit.Aborted, nil, "not overwriting your current secret")
		}
	}

	if err := s.Store.Move(ctx, from, to); err != nil {
		return exit.Error(exit.Unknown, err, "%s", err)
	}

	return nil
}
