package action

import (
	"context"
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/kpitt/gopass/internal/backend"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/urfave/cli/v2"
)

// Version prints the gopass version.
func (s *Action) Version(c *cli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)

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
