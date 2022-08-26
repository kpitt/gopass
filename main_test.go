package main

import (
	"bytes"
	"context"
	"flag"
	"os"
	"testing"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/kpitt/gopass/internal/action"
	"github.com/kpitt/gopass/internal/backend"
	"github.com/kpitt/gopass/internal/backend/crypto/gpg"
	"github.com/kpitt/gopass/internal/config"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/internal/set"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/tests/gptest"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestVersionPrinter(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	vp := makeVersionPrinter(buf, "1.0.0", "2022-08-17")
	vp(nil)
	assert.Contains(t, buf.String(), "gopass version 1.0.0 (2022-08-17)\n")
}

func TestSetupApp(t *testing.T) {
	t.Parallel()

	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	_, app := setupApp(ctx, "1.9.0", "2022-08-17")
	assert.NotNil(t, app)
}

// commandsWithError is a list of commands that return an error when
// invoked without arguments.
var commandsWithError = set.Map([]string{
	".age.identities.add",
	".age.identities.remove",
	".alias.add",
	".alias.remove",
	".alias.delete",
	".audit",
	".cat",
	".clone",
	".copy",
	".create",
	".delete",
	".edit",
	".env",
	".find",
	".fscopy",
	".fsmove",
	".generate",
	".git",
	".git.push",
	".git.pull",
	".git.status",
	".git.remote.add",
	".git.remote.remove",
	".grep",
	".history",
	".init",
	".insert",
	".link",
	".merge",
	".mounts.add",
	".mounts.remove",
	".move",
	".otp",
	".process",
	".rcs.status",
	".recipients.add",
	".recipients.remove",
	".show",
	".sum",
	".templates.edit",
	".templates.remove",
	".templates.show",
	".unclip",
})

func TestGetCommands(t *testing.T) { //nolint:paralleltest
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	buf := &bytes.Buffer{}
	color.NoColor = true

	out.Stdout = buf
	defer func() {
		out.Stdout = os.Stdout
	}()

	cfg := config.New()
	cfg.Path = u.StoreDir("")

	clipboard.Unsupported = true

	ctx := context.Background()
	ctx = ctxutil.WithAlwaysYes(ctx, true)
	ctx = ctxutil.WithInteractive(ctx, false)
	ctx = ctxutil.WithTerminal(ctx, false)
	ctx = ctxutil.WithHidden(ctx, true)
	ctx = backend.WithCryptoBackendString(ctx, "plain")

	act, err := action.New(cfg, "1.9.0")
	assert.NoError(t, err)

	app := cli.NewApp()
	fs := flag.NewFlagSet("default", flag.ContinueOnError)
	c := cli.NewContext(app, fs, nil)
	c.Context = ctx

	commands := getCommands(act, app)
	assert.Equal(t, 40, len(commands))

	prefix := ""
	testCommands(t, c, commands, prefix)
}

func testCommands(t *testing.T, c *cli.Context, commands []*cli.Command, prefix string) {
	t.Helper()

	for _, cmd := range commands {
		if len(cmd.Subcommands) > 0 {
			testCommands(t, c, cmd.Subcommands, prefix+"."+cmd.Name)
		}

		if cmd.Before != nil {
			if err := cmd.Before(c); err != nil {
				continue
			}
		}

		if cmd.BashComplete != nil {
			cmd.BashComplete(c)
		}

		if cmd.Action != nil {
			fullName := prefix + "." + cmd.Name
			if _, found := commandsWithError[fullName]; found {
				assert.Error(t, cmd.Action(c), fullName)

				continue
			}

			assert.NoError(t, cmd.Action(c), fullName)
		}
	}
}

func TestInitContext(t *testing.T) {
	t.Parallel()

	u := gptest.NewUnitTester(t)
	defer u.Remove()

	ctx := context.Background()
	cfg := config.New()

	ctx = initContext(ctx, cfg)
	assert.Equal(t, true, gpg.IsAlwaysTrust(ctx))
}
