package action

import (
	"context"
	"fmt"

	"github.com/kpitt/gopass/internal/action/exit"
	"github.com/kpitt/gopass/internal/backend"
	"github.com/kpitt/gopass/internal/backend/crypto/age"
	"github.com/kpitt/gopass/internal/backend/crypto/gpg"
	"github.com/kpitt/gopass/internal/config"
	"github.com/kpitt/gopass/internal/cui"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/debug"
	"github.com/kpitt/gopass/pkg/fsutil"
	"github.com/kpitt/gopass/pkg/termio"
	"github.com/urfave/cli/v2"
)

// IsInitialized returns an error if the store is not properly
// prepared.
func (s *Action) IsInitialized(c *cli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)
	inited, err := s.Store.IsInitialized(ctx)
	if err != nil {
		return exit.Error(exit.Unknown, err, "Failed to initialize store: %s", err)
	}

	if inited {
		debug.Log("Store is fully initialized and ready to go.\n")
		s.printReminder(ctx)
		if c.Command.Name != "sync" {
			_ = s.autoSync(ctx)
		}

		return nil
	}

	debug.Log("Store needs to be initialized.\n")
	if !ctxutil.IsInteractive(ctx) {
		return exit.Error(exit.NotInitialized, nil, "password-store is not initialized. Try '%s init'", s.Name)
	}

	out.Errorf(ctx, "No existing configuration found.")
	out.Printf(ctx, "- Please run 'gopass init'")

	return exit.Error(exit.NotInitialized, err, "not initialized")
}

// Init a new password store with a first gpg id.
func (s *Action) Init(c *cli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)
	path := c.String("path")
	alias := c.String("store")
	remoteUrl := c.String("remote")

	ctx = initParseContext(ctx, c)
	out.Printf(ctx, "Initializing a new password store:\n")

	if name := termio.DetectName(c.Context, c); name != "" {
		ctx = ctxutil.WithUsername(ctx, name)
	}

	if email := termio.DetectEmail(c.Context, c); email != "" {
		ctx = ctxutil.WithEmail(ctx, email)
	}

	// age: only native keys
	// "[ssh] types should only be used for compatibility with existing keys,
	// and native X25519 keys should be preferred otherwise."
	// https://pkg.go.dev/filippo.io/age@v1.0.0/agessh#pkg-overview.
	ctx = age.WithOnlyNative(ctx, true)
	// gpg: only trusted keys
	// only list "usable" / properly trused and signed GPG keys by requesting
	// always trust is false. Ignored for other backends. See
	// https://www.gnupg.org/gph/en/manual/r1554.html.
	ctx = gpg.WithAlwaysTrust(ctx, false)

	inited, err := s.Store.IsInitialized(ctx)
	if err != nil {
		return exit.Error(exit.Unknown, err, "Failed to initialized store: %s", err)
	}

	if inited {
		out.Errorf(ctx, "Store is already initialized!")
	}

	crypto := s.getCryptoFor(ctx, alias)
	if crypto == nil {
		return fmt.Errorf("cannot continue without crypto")
	}
	debug.Log("Crypto Backend initialized as: %s", crypto.Name())

	if err := s.initCheckPrivateKeys(ctx, crypto); err != nil {
		return fmt.Errorf("failed to check private keys: %w", err)
	}

	if err := s.init(ctx, alias, path, remoteUrl, c.Args().Slice()...); err != nil {
		return exit.Error(exit.Unknown, err, "Failed to initialize store: %s", err)
	}

	return nil
}

func initParseContext(ctx context.Context, c *cli.Context) context.Context {
	if c.IsSet("crypto") {
		ctx = backend.WithCryptoBackendString(ctx, c.String("crypto"))
	}

	if c.IsSet("storage") {
		ctx = backend.WithStorageBackendString(ctx, c.String("storage"))
	}

	if !backend.HasCryptoBackend(ctx) {
		debug.Log("Using default Crypto Backend (GPGCLI)")
		ctx = backend.WithCryptoBackend(ctx, backend.GPGCLI)
	}

	if !backend.HasStorageBackend(ctx) {
		debug.Log("Using default storage backend (GitFS)")
		ctx = backend.WithStorageBackend(ctx, backend.GitFS)
	}

	return ctx
}

func (s *Action) init(ctx context.Context, alias, path string, remoteUrl string, keys ...string) error {
	if path == "" {
		if alias != "" {
			path = config.PwStoreDir(alias)
		} else {
			path = s.Store.Path()
		}
	}
	path = fsutil.CleanPath(path)

	debug.Log("Initializing Store %q in %q for %+v", alias, path, keys)

	out.Printf(ctx, "- Searching for usable private keys...")
	debug.Log("Checking private keys for: %+v", keys)
	crypto := s.getCryptoFor(ctx, alias)

	// private key selection doesn't matter for plain. save one question.
	// TODO should ask the backend
	if crypto.Name() == "plain" {
		keys, _ = crypto.ListIdentities(ctx)
	}

	if len(keys) < 1 {
		out.Notice(ctx, "Hint: Use 'gopass init <subkey> to use subkeys!'")
		nk, err := cui.AskForPrivateKey(ctx, crypto, "? Please select a private key for encrypting secrets:")
		if err != nil {
			return fmt.Errorf("failed to read user input: %w", err)
		}
		keys = []string{nk}
	}

	debug.Log("Initializing sub store - Alias: %q - Path: %q - Keys: %+v", alias, path, keys)
	if err := s.Store.Init(ctx, alias, path, keys...); err != nil {
		return fmt.Errorf("failed to init store %q at %q: %w", alias, path, err)
	}

	if alias != "" && path != "" {
		debug.Log("Mounting sub store %q -> %q", alias, path)
		if err := s.Store.AddMount(ctx, alias, path); err != nil {
			return fmt.Errorf("failed to add mount %q: %w", alias, err)
		}
	}

	be := backend.GetStorageBackend(ctx)
	if be == backend.GitFS {
		debug.Log("Initializing git repository...")
		if err := s.rcsInit(ctx, alias, ctxutil.GetUsername(ctx), ctxutil.GetEmail(ctx)); err != nil {
			debug.Log("Stacktrace: %+v\n", err)
			out.Errorf(ctx, "✗ Failed to initialize git repository: %s", err)
		}
		debug.Log("Git initialized as %s", s.Store.Storage(ctx, alias).Name())
	} else {
		debug.Log("not initializing git backend")
	}

	// write config.
	debug.Log("Writing configuration to %q", s.cfg.ConfigPath)
	if err := s.cfg.Save(); err != nil {
		return exit.Error(exit.Config, err, "failed to write config: %s", err)
	}

	out.Printf(ctx, "✓ Password store %s initialized for:", path)
	s.printRecipients(ctx, alias)

	if be == backend.GitFS && remoteUrl != "" {
		debug.Log("configuring git remote: %q", remoteUrl)
		out.Printf(ctx, "Configuring git remote...")
		if err := s.initSetupGitRemote(ctx, alias, remoteUrl); err != nil {
			return fmt.Errorf("failed to setup git remote: %w", err)
		}
	}

	return nil
}

func (s *Action) initSetupGitRemote(ctx context.Context, alias, remoteUrl string) error {
	// omit RCS output.
	ctx = ctxutil.WithHidden(ctx, true)
	if err := s.Store.RCSAddRemote(ctx, alias, "origin", remoteUrl); err != nil {
		return fmt.Errorf("failed to add git remote: %w", err)
	}
	// initial pull, in case the remote is non-empty.
	if err := s.Store.RCSPull(ctx, alias, "origin", ""); err != nil {
		debug.Log("Initial git pull failed: %s", err)
	}
	if err := s.Store.RCSPush(ctx, alias, "origin", ""); err != nil {
		return fmt.Errorf("failed to push to git remote: %w", err)
	}

	return nil
}

func (s *Action) printRecipients(ctx context.Context, alias string) {
	crypto := s.Store.Crypto(ctx, alias)
	for _, recipient := range s.Store.ListRecipients(ctx, alias) {
		r := "0x" + recipient
		if kl, err := crypto.FindRecipients(ctx, recipient); err == nil && len(kl) > 0 {
			r = crypto.FormatKey(ctx, kl[0], "")
		}
		out.Printf(ctx, "- "+r)
	}
}

func (s *Action) getCryptoFor(ctx context.Context, name string) backend.Crypto {
	return s.Store.Crypto(ctx, name)
}
