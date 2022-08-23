package action

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kpitt/gopass/internal/action/exit"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/internal/store"
	"github.com/kpitt/gopass/pkg/clipboard"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/debug"
	"github.com/kpitt/gopass/pkg/otp"
	"github.com/kpitt/gopass/pkg/termio"
	"github.com/mattn/go-tty"
	potp "github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
	"github.com/urfave/cli/v2"
)

// OTP implements OTP token handling for TOTP and HOTP.
func (s *Action) OTP(c *cli.Context) error {
	ctx := ctxutil.WithGlobalFlags(c)
	name := c.Args().First()
	if name == "" {
		return exit.Error(exit.Usage, nil, "Usage: %s otp <NAME>", s.Name)
	}

	qrf := c.String("qr")
	clip := c.Bool("clip")
	pw := c.Bool("password")

	return s.otp(ctx, name, qrf, clip, pw, true)
}

func tickingBar(ctx context.Context, expiresAt time.Time, bar *termio.ProgressBar) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for tt := range ticker.C {
		select {
		case <-ctx.Done():
			return // returning not to leak the goroutine.
		default:
			// we don't want to block if not cancelled.
		}
		if tt.After(expiresAt) {
			return
		}
		bar.Inc()
	}
}

func waitForKeyPress(ctx context.Context, cancel context.CancelFunc) {
	tty, err := tty.Open()
	if err != nil {
		out.Errorf(ctx, "Unexpected error opening tty: %v", err)
		cancel()
	}

	defer func() {
		_ = tty.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return // returning not to leak the goroutine.
		default:
		}

		r, err := tty.ReadRune()
		if err != nil {
			out.Errorf(ctx, "Unexpected error opening tty: %v", err)
		}

		if r == 'q' || r == 'x' || err != nil {
			cancel()

			return
		}
	}
}

//nolint:cyclop
func (s *Action) otp(ctx context.Context, name, qrf string, clip, pw, recurse bool) error {
	sec, err := s.Store.Get(ctx, name)
	if err != nil {
		return s.otpHandleError(ctx, name, qrf, clip, pw, recurse, err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	skip := ctxutil.IsHidden(ctx) || pw || qrf != "" || !ctxutil.IsTerminal(ctx) || !ctxutil.IsInteractive(ctx) || clip
	if !skip {
		// let us monitor key presses for cancellation:.
		out.Warningf(ctx, "%s", "[q] to stop. -o flag to avoid.")
		go waitForKeyPress(ctx, cancel)
	}

	// only used for the HOTP case as a fallback
	var counter uint64 = 1
	if sv, found := sec.Get("counter"); found && sv != "" {
		if iv, err := strconv.ParseUint(sv, 10, 64); iv != 0 && err == nil {
			counter = iv
		}
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		two, err := otp.Calculate(name, sec)
		if err != nil {
			return exit.Error(exit.Unknown, err, "No OTP entry found for %s: %s", name, err)
		}

		var token string
		switch two.Type() {
		case "totp":
			token, err = totp.GenerateCodeCustom(two.Secret(), time.Now(), totp.ValidateOpts{
				Period:    uint(two.Period()),
				Skew:      1,
				Digits:    parseDigits(two.URL()),
				Algorithm: parseAlgorithm(two.URL()),
			})
			if err != nil {
				return exit.Error(exit.Unknown, err, "Failed to compute OTP token for %s: %s", name, err)
			}
		case "hotp":
			token, err = hotp.GenerateCodeCustom(two.Secret(), counter, hotp.ValidateOpts{
				Digits:    parseDigits(two.URL()),
				Algorithm: parseAlgorithm(two.URL()),
			})
			if err != nil {
				return exit.Error(exit.Unknown, err, "Failed to compute OTP token for %s: %s", name, err)
			}
			counter++
			_ = sec.Set("counter", strconv.Itoa(int(counter)))
			if err := s.Store.Set(ctx, name, sec); err != nil {
				out.Errorf(ctx, "Failed to persist counter value: %s", err)
			}
			debug.Log("Saved counter as %d", counter)
		}

		now := time.Now()
		expiresAt := now.Add(time.Duration(two.Period()) * time.Second).Truncate(time.Duration(two.Period()) * time.Second)
		secondsLeft := int(time.Until(expiresAt).Seconds())
		bar := termio.NewProgressBar(token, int64(secondsLeft))
		bar.Hidden = skip

		debug.Log("OTP period: %ds", two.Period())

		if clip {
			if err := clipboard.CopyTo(ctx, fmt.Sprintf("token for %s", name), []byte(token), s.cfg.ClipTimeout); err != nil {
				return exit.Error(exit.IO, err, "failed to copy to clipboard: %s", err)
			}

			return nil
		}

		// check if we are in "password only" or in "qr code" mode or being redirected to a pipe.
		if pw || qrf != "" || !ctxutil.IsTerminal(ctx) {
			out.Printf(ctx, "%s", token)
			cancel()
		} else { // if not then we want to print a progress bar with the expiry time.
			if bar.Hidden {
				cancel()
			} else {
				bar.Set(0)
				go tickingBar(ctx, expiresAt, bar)
			}
		}

		if qrf != "" {
			return otp.WriteQRFile(two, qrf)
		}

		// let us wait until next OTP code:.
		for {
			select {
			case <-ctx.Done():
				bar.Done()

				return nil
			default:
				time.Sleep(time.Millisecond * 500)
			}
			if time.Now().After(expiresAt) {
				bar.Done()

				break
			}
		}
	}
}

func (s *Action) otpHandleError(ctx context.Context, name, qrf string, clip, pw, recurse bool, err error) error {
	if !errors.Is(err, store.ErrNotFound) || !recurse || !ctxutil.IsTerminal(ctx) {
		return exit.Error(exit.Unknown, err, "failed to retrieve secret %q: %s", name, err)
	}

	out.Printf(ctx, "Entry %q not found. Starting search...", name)
	cb := func(ctx context.Context, c *cli.Context, name string, recurse bool) error {
		return s.otp(ctx, name, qrf, clip, pw, false)
	}
	if err := s.find(ctx, nil, name, cb, false); err != nil {
		return exit.Error(exit.NotFound, err, "%s", err)
	}

	return nil
}

// parseDigits and parseAlgorithm can be replaced if https://github.com/pquerna/otp/pull/74 is merged.
func parseDigits(ku string) potp.Digits {
	u, err := url.Parse(ku)
	if err != nil {
		debug.Log("Failed to parse key URL: %s", err)

		// return the most common value
		return potp.DigitsSix
	}

	q := u.Query()
	iv, err := strconv.ParseUint(q.Get("digits"), 10, 64)
	if err != nil {
		debug.Log("Failed to parse digits param: %s", err)

		// return the most common value
		return potp.DigitsSix
	}

	switch iv {
	case 6:
		return potp.DigitsSix
	case 8:
		return potp.DigitsEight
	default:
		debug.Log("Unsupported digits value: %d", iv)

		// return the most common value
		return potp.DigitsSix
	}
}

func parseAlgorithm(ku string) potp.Algorithm {
	u, err := url.Parse(ku)
	if err != nil {
		debug.Log("Failed to parse key URL: %s", err)

		// return the most common value
		return potp.AlgorithmSHA1
	}

	q := u.Query()
	a := strings.ToLower(q.Get("algorithm"))
	switch a {
	case "md5":
		return potp.AlgorithmMD5
	case "sha256":
		return potp.AlgorithmSHA256
	case "sha512":
		return potp.AlgorithmSHA512
	default:
		return potp.AlgorithmSHA1
	}
}
