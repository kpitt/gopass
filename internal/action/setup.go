package action

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kpitt/gopass/internal/backend"
	"github.com/kpitt/gopass/internal/backend/crypto/gpg"
	gpgcli "github.com/kpitt/gopass/internal/backend/crypto/gpg/cli"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/debug"
	"github.com/kpitt/gopass/pkg/pwgen/xkcdgen"
	"github.com/kpitt/gopass/pkg/termio"
)

func (s *Action) initCheckPrivateKeys(ctx context.Context, crypto backend.Crypto) error {
	// check for existing GPG/Age keypairs (private/secret keys). We need at least
	// one useable key pair. If none exists try to create one.
	if !s.initHasUseablePrivateKeys(ctx, crypto) {
		out.Printf(ctx, "! No useable cryptographic keys. Generating new key pair")
		if crypto.Name() == "gpgcli" {
			out.Printf(ctx, "! Key generation may take up to a few minutes")
		}
		if err := s.initGenerateIdentity(ctx, crypto, ctxutil.GetUsername(ctx), ctxutil.GetEmail(ctx)); err != nil {
			return fmt.Errorf("failed to create new private key: %w", err)
		}
		out.Printf(ctx, "✓ Cryptographic keys generated")
	}

	debug.Log("We have useable private keys")

	return nil
}

func (s *Action) initGenerateIdentity(ctx context.Context, crypto backend.Crypto, name, email string) error {
	out.Printf(ctx, "- Creating cryptographic key pair (%s)...", crypto.Name())

	if crypto.Name() == gpgcli.Name {
		var err error

		out.Printf(ctx, "- Gathering information for the %s key pair...", crypto.Name())
		name, err = termio.AskForString(ctx, "? What is your name?", name)
		if err != nil {
			return err
		}

		email, err = termio.AskForString(ctx, "? What is your email?", email)
		if err != nil {
			return err
		}
	}

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

	if crypto.Name() == "gpgcli" {
		// Note: This issue shouldn't matter much past Linux Kernel 5.6,
		// eventually we might want to remove this notice. Only applies to
		// GPG.
		out.Printf(ctx, "! This can take a long time. If you get impatient see https://gopass.pittcrew.org/entropy")
		if want, err := termio.AskForBool(ctx, "Continue?", true); err != nil || !want {
			return fmt.Errorf("user aborted: %w", err)
		}
	}

	if err := crypto.GenerateIdentity(ctx, name, email, passphrase); err != nil {
		return fmt.Errorf("failed to create new private key: %w", err)
	}

	out.OKf(ctx, "Key pair generated")

	if pwGenerated {
		out.Printf(ctx, color.MagentaString("Passphrase: ")+passphrase)
		out.Noticef(ctx, "You need to remember this very well!")
	}

	out.Notice(ctx, "! We need to unlock your newly created private key now! Please enter the passphrase you just generated.")

	// avoid the gpg cache or we won't find the newly created key
	kl, err := crypto.ListIdentities(gpg.WithUseCache(ctx, false))
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
