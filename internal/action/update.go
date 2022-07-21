package action

import (
	"fmt"
	"runtime"

	"github.com/kpitt/gopass/internal/action/exit"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/internal/updater"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/urfave/cli/v2"
)

// Update will start the interactive update assistant.
func (s *Action) Update(c *cli.Context) error {
	_ = s.rem.Reset("update")

	ctx := ctxutil.WithGlobalFlags(c)

	if s.version.String() == "0.0.0+HEAD" {
		out.Errorf(ctx, "Can not check version against HEAD")

		return nil
	}

	if runtime.GOOS == "windows" {
		return fmt.Errorf("gopass update is not supported on windows (#1722)")
	}

	out.Printf(ctx, "⚒ Checking for available updates ...")
	if err := updater.Update(ctx, s.version); err != nil {
		return exit.Error(exit.Unknown, err, "Failed to update gopass: %s", err)
	}

	out.OKf(ctx, "gopass is up to date")

	return nil
}
