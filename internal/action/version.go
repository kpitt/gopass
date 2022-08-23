package action

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/fatih/color"
	"github.com/kpitt/gopass/internal/action/exit"
	"github.com/kpitt/gopass/internal/backend"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/internal/updater"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/debug"
	"github.com/kpitt/gopass/pkg/protect"
	"github.com/urfave/cli/v2"
)

// Version prints the gopass version.
func (s *Action) Version(c *cli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)
	version := make(chan string, 1)
	go s.checkVersion(ctx, version)

	_ = s.IsInitialized(c)

	cli.VersionPrinter(c)
	fmt.Fprintln(stdout)

	cryptoVer := versionInfo(ctx, s.Store.Crypto(ctx, ""))
	storageVer := versionInfo(ctx, s.Store.Storage(ctx, ""))

	tpl := "%-10s - %10s - %10s\n"
	fmt.Fprintf(stdout, tpl, "<root>", cryptoVer, storageVer)

	// report all used crypto, sync and fs backends.
	for _, mp := range s.Store.MountPoints() {
		cv := versionInfo(ctx, s.Store.Crypto(ctx, mp))
		sv := versionInfo(ctx, s.Store.Storage(ctx, mp))

		if cv != cryptoVer || sv != storageVer {
			fmt.Fprintf(stdout, tpl, mp, cv, sv)
		}
	}

	fmt.Fprintln(stdout)
	fmt.Fprintf(stdout, "Available Crypto Backends: %s\n", strings.Join(backend.CryptoRegistry.BackendNames(), ", "))
	fmt.Fprintf(stdout, "Available Storage Backends: %s\n", strings.Join(backend.StorageRegistry.BackendNames(), ", "))

	select {
	case vi := <-version:
		if vi != "" {
			fmt.Fprintln(stdout, vi)
		}
	case <-time.After(2 * time.Second):
		out.Errorf(ctx, "Version check timed out")
	case <-ctx.Done():
		return exit.Error(exit.Aborted, nil, "user aborted")
	}

	return nil
}

type versioner interface {
	Name() string
	Version(context.Context) semver.Version
}

func versionInfo(ctx context.Context, v versioner) string {
	if v == nil {
		return "<none>"
	}

	return fmt.Sprintf("%s %s", v.Name(), v.Version(ctx))
}

func (s *Action) getSemver() (semver.Version, error) {
	version := strings.TrimPrefix(s.version, "v")
	parts := strings.SplitN(version, "-", 3)
	if len(parts) != 3 {
		// doesn't look like a "git describe" version, so parse as-is
		return semver.Parse(version)
	}

	baseVer, err := semver.Parse(parts[0])
	if err != nil {
		return baseVer, err
	}

	// Treat "git describe" dev version as a pre-release of the next patch
	baseVer.Patch++
	version = fmt.Sprintf("%s-DEV.%s+%s", baseVer.String(), parts[1], parts[2])

	return semver.Parse(version)
}

func (s *Action) checkVersion(ctx context.Context, u chan string) {
	msg := ""
	defer func() {
		u <- msg
	}()

	if disabled := os.Getenv("CHECKPOINT_DISABLE"); disabled != "" {
		debug.Log("remote version check disabled by CHECKPOINT_DISABLE")

		return
	}

	// force checking for updates, mainly for testing.
	force := os.Getenv("GOPASS_FORCE_CHECK") != ""

	if !force && strings.ContainsRune(s.version, '-') {
		// chan not check version against HEAD.
		debug.Log("remote version check disabled for dev version")

		return
	}

	if !force && protect.ProtectEnabled {
		// chan not check version
		// against pledge(2)'d OpenBSD.
		debug.Log("remote version check disabled for pledge(2)'d version")

		return
	}

	sv, err := s.getSemver()
	if err != nil {
		msg = color.RedString("\nError parsing current version: %s", err)
		return
	}

	r, err := updater.FetchLatestRelease(ctx)
	if err != nil {
		msg = color.RedString("\nError checking latest version: %s", err)
		return
	}

	if sv.GTE(r.Version) {
		_ = s.rem.Reset("update")
		debug.Log("gopass is up-to-date (local: %q, GitHub: %q)", s.version, r.Version)

		return
	}

	notice := fmt.Sprintf("\nYour version (%s) of gopass is out of date!\nThe latest version is %s.\n", s.version, r.Version.String())
	notice += "You can update by downloading from https://www.gopass.pw/#install"
	if err := updater.IsUpdateable(ctx); err == nil {
		notice += " by running 'gopass update'"
	}
	notice += " or via your package manager"
	msg = color.YellowString(notice)
}
