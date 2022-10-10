package action

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kpitt/gopass/internal/backend"
	gpgcli "github.com/kpitt/gopass/internal/backend/crypto/gpg/cli"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/pkg/debug"
	"github.com/kpitt/gopass/pkg/pwgen/xkcdgen"
	"github.com/kpitt/gopass/pkg/termio"
)

func (s *Action) initCheckPrivateKeys(ctx context.Context, crypto backend.Crypto) error {
	// check for existing GPG/Age keypairs (private/secret keys). We need at least
	// one useable key pair. If none exists, try to create one if using `age` encryption.
	if !s.initHasUseablePrivateKeys(ctx, crypto) {
		if crypto.Name() == gpgcli.Name {
			debug.Log("Setup Error: no useable GPG private keys found")
			out.Error(ctx, "No useable cryptographic keys")

			// TODO: Handle initGenerateIdentity not supported
			return backend.ErrNotFound
		}

		out.Printf(ctx, "! No useable cryptographic keys. Generating a new key pair.")
		if err := s.initGenerateIdentity(ctx, crypto); err != nil {
			return fmt.Errorf("failed to create new private key: %w", err)
		}
		out.Printf(ctx, "✓ Cryptographic keys generated")
	}

	debug.Log("We have useable private keys")

	return nil
}

func (s *Action) initGenerateIdentity(ctx context.Context, crypto backend.Crypto) error {
	passphrase := xkcdgen.Random()
	pwGenerated := true
	want, err := termio.AskForBool(ctx, "! Do you want to enter a passphrase? (otherwise we generate one for you)", false)
	if err != nil {
		return err
	}
	if want {
		pwGenerated = false
		sv, err := termio.AskForPassword(ctx, "passphrase for your new keypair", true)
		if err != nil {
			return fmt.Errorf("failed to read passphrase: %w", err)
		}
		passphrase = sv
	}

	if err := crypto.GenerateIdentity(ctx, passphrase); err != nil {
		return fmt.Errorf("failed to create new private key: %w", err)
	}

	out.OKf(ctx, "Key pair generated")

	if pwGenerated {
		out.Printf(ctx, color.MagentaString("Passphrase: ")+passphrase)
		out.Noticef(ctx, "You need to remember this very well!")
	}

	out.Notice(ctx, "! We need to unlock your newly created private key now! Please enter the passphrase you just generated.")

	kl, err := crypto.ListIdentities(ctx)
	if err != nil {
		return fmt.Errorf("failed to list private keys: %w", err)
	}

	if len(kl) > 1 {
		out.Notice(ctx, "More than one private key detected. Make sure to use the correct one!")
	}

	if len(kl) < 1 {
		return fmt.Errorf("failed to create a usable key pair")
	}

	// we can export the generated key to the current directory for convenience.
	if err := s.initExportPublicKey(ctx, crypto, kl[0]); err != nil {
		return err
	}
	out.OKf(ctx, "Key pair validated")

	return nil
}

type keyExporter interface {
	ExportPublicKey(ctx context.Context, id string) ([]byte, error)
}

func (s *Action) initExportPublicKey(ctx context.Context, crypto backend.Crypto, key string) error {
	exp, ok := crypto.(keyExporter)
	if !ok {
		debug.Log("crypto backend %T cannot export public keys", crypto)

		return nil
	}

	fn := key + ".pub.key"
	want, err := termio.AskForBool(ctx, fmt.Sprintf("Do you want to export your public key to %q?", fn), false)
	if err != nil {
		return err
	}

	if !want {
		return nil
	}

	pk, err := exp.ExportPublicKey(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to export public key: %w", err)
	}

	if err := os.WriteFile(fn, pk, 0o6444); err != nil {
		out.Errorf(ctx, "✗ Failed to export public key %q: %q", fn, err)

		return err
	}
	out.Printf(ctx, "✓ Public key exported to %q", fn)

	return nil
}

func (s *Action) initHasUseablePrivateKeys(ctx context.Context, crypto backend.Crypto) bool {
	debug.Log("checking for existing, usable identities / private keys for %s", crypto.Name())
	kl, err := crypto.ListIdentities(ctx)
	if err != nil {
		return false
	}

	debug.Log("available private keys: %q for %s", kl, crypto.Name())

	return len(kl) > 0
}
