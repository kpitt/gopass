package action

import (
	"context"
	"os"

	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/pkg/ctxutil"
)

func (s *Action) printReminder(ctx context.Context) {
	if !ctxutil.IsInteractive(ctx) {
		return
	}

	if !ctxutil.IsTerminal(ctx) {
		return
	}

	if sv := os.Getenv("GOPASS_NO_REMINDER"); sv != "" {
		return
	}

	// Note: We only want to print one reminder per day (at most).
	// So we intentionally return after printing one, leaving the others
	// for the following days.
	if s.rem.Overdue("fsck") {
		out.Notice(ctx, "You haven't run 'gopass fsck' in a while.")

		return
	}

	if s.rem.Overdue("audit") {
		out.Notice(ctx, "You haven't run 'gopass audit' in a while.")

		return
	}
}
