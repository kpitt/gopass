package action

import (
	"fmt"
	"runtime"
	"strings"

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

	if strings.ContainsRune(s.version, '-') {
		out.Errorf(ctx, "Cannot check pre-release version")
		return nil
	}

	if runtime.GOOS == "windows" {
		return fmt.Errorf("gopass update is not supported on windows (#1722)")
	}

	sv, err := s.getSemver()
	if err != nil {
		return fmt.Errorf("could not parse current gopass version")
	}

	out.Printf(ctx, "- Checking for available updates...")
	if err := updater.Update(ctx, sv); err != nil {
		return exit.Error(exit.Unknown, err, "Failed to update gopass: %s", err)
	}

	out.OKf(ctx, "gopass is up to date")

	return nil
}
