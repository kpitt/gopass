package action

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/kpitt/gopass/internal/backend"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/internal/store"
	"github.com/kpitt/gopass/internal/store/leaf"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/debug"
	"github.com/urfave/cli/v2"
)

var autosyncIntervalDays = 3

func init() {
	sv := os.Getenv("GOPASS_AUTOSYNC_INTERVAL")
	if sv == "" {
		return
	}

	iv, err := strconv.Atoi(sv)
	if err != nil {
		return
	}

	autosyncIntervalDays = iv
}

// Sync all stores with their remotes.
func (s *Action) Sync(c *cli.Context) error {
	return s.sync(ctxutil.WithGlobalFlags(c), c.String("store"))
}

func (s *Action) autoSync(ctx context.Context) error {
	if !ctxutil.IsInteractive(ctx) {
		return nil
	}

	if !ctxutil.IsTerminal(ctx) {
		return nil
	}

	if sv := os.Getenv("GOPASS_NO_AUTOSYNC"); sv != "" {
		return nil
	}

	ls := s.rem.LastSeen("autosync")
	debug.Log("autosync - last seen: %s", ls)
	if time.Since(ls) > time.Duration(autosyncIntervalDays)*24*time.Hour {
		return s.sync(ctx, "")
	}

	return nil
}

func (s *Action) sync(ctx context.Context, store string) error {
	mps := s.Store.MountPoints()
	mps = append([]string{""}, mps...)

	// sync all stores (root and all mounted sub stores).
	for _, mp := range mps {
		if store != "" {
			if store != "root" && mp != store {
				continue
			}
			if store == "root" && mp != "" {
				continue
			}
		}

		_ = s.syncMount(ctx, mp)
	}
	out.OKf(ctx, "All done")

	// If we are sync'ing all stores, and sync succeeds, then reset the auto-sync interval.
	if store == "" {
		_ = s.rem.Reset("autosync")
	}

	return nil
}

// syncMount syncs a single mount.
func (s *Action) syncMount(ctx context.Context, mp string) error {
	ctxno := out.WithNewline(ctx, false)
	name := mp
	if mp == "" {
		name = "<root>"
	}

	sub, err := s.Store.GetSubStore(mp)
	if err != nil {
		out.Errorf(ctx, "Failed to get sub store %q: %s", name, err)

		return fmt.Errorf("failed to get sub stores (%w)", err)
	}

	if sub == nil {
		out.Errorf(ctx, "Failed to get sub stores '%s: nil'", name)

		return fmt.Errorf("failed to get sub stores (nil)")
	}

	syncMsg := fmt.Sprintf("Synchronizing %s store", color.CyanString(name))
	ctx = ctxutil.WithSpinner(ctx, syncMsg)

	err = sub.Storage().Push(ctx, "", "")
	switch {
	case err == nil:
		debug.Log("Push succeeded")
	case errors.Is(err, store.ErrGitNoRemote):
		out.Noticef(ctx, "Skipped %q store (no remote)", name)
		debug.Log("Failed to push %q to its remote: %s", name, err)

		return err
	case errors.Is(err, backend.ErrNotSupported):
		break
	case errors.Is(err, store.ErrGitNotInit):
		break
	default: // any other error
		out.Errorf(ctx, "Failed to push %q to its remote: %s", name, err)

		return err
	}

	debug.Log("Syncing Mount %s. Exportkeys: %t", mp, ctxutil.IsExportKeys(ctx))
	if err := syncImportKeys(ctxno, sub, name); err != nil {
		return err
	}
	if ctxutil.IsExportKeys(ctx) {
		if err := syncExportKeys(ctxno, sub, name); err != nil {
			return err
		}
	}
	out.OKf(ctx, "%s %s", syncMsg, color.GreenString("[OK]"))

	return nil
}

func syncImportKeys(ctx context.Context, sub *leaf.Store, name string) error {
	// import keys.
	if err := sub.ImportMissingPublicKeys(ctx); err != nil {
		out.Errorf(ctx, "Failed to import missing public keys for %q: %s", name, err)

		return err
	}

	return nil
}

func syncExportKeys(ctx context.Context, sub *leaf.Store, name string) error {
	// export keys.
	rs, err := sub.GetRecipients(ctx, "")
	if err != nil {
		out.Errorf(ctx, "Failed to load recipients for %q: %s", name, err)

		return err
	}
	exported, err := sub.ExportMissingPublicKeys(ctx, rs)
	if err != nil {
		out.Errorf(ctx, "Failed to export missing public keys for %q: %s", name, err)

		return err
	}

	// only run second push if we did export any keys.
	if !exported {
		return nil
	}

	if err := sub.Storage().Push(ctx, "", ""); err != nil {
		out.Errorf(ctx, "Failed to push %q to its remote: %s", name, err)

		return err
	}

	return nil
}
