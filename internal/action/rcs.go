package action

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kpitt/gopass/internal/action/exit"
	"github.com/kpitt/gopass/internal/backend"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/debug"
	"github.com/kpitt/gopass/pkg/termio"
	"github.com/urfave/cli/v2"
)

// RCSInit initializes a git repo including basic configuration.
func (s *Action) RCSInit(c *cli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)
	store := c.String("store")
	un := termio.DetectName(c.Context, c)
	ue := termio.DetectEmail(c.Context, c)

	// default to git.
	if !backend.HasStorageBackend(ctx) {
		ctx = backend.WithStorageBackend(ctx, backend.GitFS)
	}

	if err := s.rcsInit(ctx, store, un, ue); err != nil {
		return exit.Error(exit.Git, err, "failed to initialize git: %s", err)
	}

	return nil
}

func (s *Action) rcsInit(ctx context.Context, store, un, ue string) error {
	userName, userEmail := s.getUserData(ctx, store, un, ue)
	if err := s.Store.RCSInit(ctx, store, userName, userEmail); err != nil {
		if errors.Is(err, backend.ErrNotSupported) {
			debug.Log("RCSInit not supported for storage backend in %q", store)

			return nil
		}

		if gtv := os.Getenv("GPG_TTY"); gtv == "" {
			out.Printf(ctx, "Git initialization failed. You may want to try to 'export GPG_TTY=$(tty)' and start over.")
		}

		return fmt.Errorf("failed to run git init: %w", err)
	}

	out.Printf(ctx, "Initialized git repository for %q <%s>...", userName, userEmail)

	return nil
}

func (s *Action) getUserData(ctx context.Context, store, name, email string) (string, string) {
	if name != "" && email != "" {
		debug.Log("Username: %s, Email: %s (provided)", name, email)

		return name, email
	}

	if name == "" {
		defaultName := termio.DetectName(ctx, nil)

		var err error
		name, err = termio.AskForString(ctx, color.CyanString("Please enter a user name for password store git config"), defaultName)
		if err != nil {
			out.Errorf(ctx, "Failed to ask for user input: %s", err)
		}
	}

	if email == "" {
		defaultEmail := termio.DetectEmail(ctx, nil)

		var err error
		email, err = termio.AskForString(ctx, color.CyanString("Please enter an email address for password store git config"), defaultEmail)
		if err != nil {
			out.Errorf(ctx, "Failed to ask for user input: %s", err)
		}
	}

	debug.Log("Username: %s, Email: %s (entered)", name, email)

	return name, email
}
