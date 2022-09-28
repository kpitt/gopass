package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigs(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		cfg  string
		want *Config
	}{
		{
			name: "1.14.4",
			cfg: `autoclip: true
autoimport: false
cliptimeout: 45
exportkeys: true
nopager: false
notifications: true
parsing: true
path: /home/johndoe/.password-store
safecontent: false
mounts:
  foo/sub: /home/johndoe/.password-store-foo-sub
  work: /home/johndoe/.password-store-work`,
			want: &Config{
				AutoClip:    true,
				AutoImport:  false,
				ClipTimeout: 45,
				ExportKeys:  true,
				NoPager:     false,
				Parsing:     true,
				Path:        "/home/johndoe/.password-store",
				Mounts: map[string]string{
					"foo/sub": "/home/johndoe/.password-store-foo-sub",
					"work":    "/home/johndoe/.password-store-work",
				},
			},
		}, {
			name: "N+1",
			cfg: `autoclip: true
autoimport: false
cliptimeout: 45
exportkeys: true
nopager: false
foo: bar
path: /home/johndoe/.password-store
mounts:
  foo/sub: /home/johndoe/.password-store-foo-sub
  work: /home/johndoe/.password-store-work`,
			want: &Config{
				AutoClip:    true,
				AutoImport:  false,
				ClipTimeout: 45,
				ExportKeys:  true,
				NoPager:     false,
				Parsing:     true,
				Path:        "/home/johndoe/.password-store",
				Mounts: map[string]string{
					"foo/sub": "/home/johndoe/.password-store-foo-sub",
					"work":    "/home/johndoe/.password-store-work",
				},
				XXX: map[string]any{"foo": string("bar")},
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := decode([]byte(tc.cfg), true)
			require.NoError(t, err)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("decode(%s) mismatch for:\n%s\n(-want +got):\n%s", tc.name, tc.cfg, diff)
			}
		})
	}
}

const testConfig = `autoclip: true
autoimport: true
cliptimeout: 5
exportkeys: true
nopager: true
notifications: true
parsing: true
path: /home/johndoe/.password-store
safecontent: true
mounts:
  foo/sub: /home/johndoe/.password-store-foo-sub
  work: /home/johndoe/.password-store-work`

func TestLoad(t *testing.T) { //nolint:paralleltest
	td := os.TempDir()
	gcfg := filepath.Join(td, ".gopass.yml")
	_ = os.Remove(gcfg)
	assert.NoError(t, os.Setenv("GOPASS_CONFIG", gcfg))
	assert.NoError(t, os.Setenv("GOPASS_HOMEDIR", td))

	require.NoError(t, os.WriteFile(gcfg, []byte(testConfig), 0o600))

	cfg := Load()
	assert.True(t, cfg.NoPager)
}

func TestLoadError(t *testing.T) { //nolint:paralleltest
	gcfg := filepath.Join(os.TempDir(), ".gopass-err.yml")
	assert.NoError(t, os.Setenv("GOPASS_CONFIG", gcfg))

	_ = os.Remove(gcfg)

	capture(t, func() error {
		_, err := load(gcfg, false)
		if err == nil {
			return fmt.Errorf("should fail")
		}

		return nil
	})

	_ = os.Remove(gcfg)
	cfg, err := load(gcfg, false)
	assert.Error(t, err)

	td, err := os.MkdirTemp("", "gopass-")
	require.NoError(t, err)

	defer func() {
		_ = os.RemoveAll(td)
	}()

	gcfg = filepath.Join(td, "foo", ".gopass.yml")
	assert.NoError(t, os.Setenv("GOPASS_CONFIG", gcfg))
	assert.NoError(t, cfg.Save())
}

func TestDecodeError(t *testing.T) { //nolint:paralleltest
	gcfg := filepath.Join(os.TempDir(), ".gopass-err2.yml")
	assert.NoError(t, os.Setenv("GOPASS_CONFIG", gcfg))

	_ = os.Remove(gcfg)
	require.NoError(t, os.WriteFile(gcfg, []byte(testConfig+"\nfoobar: zab\n"), 0o600))

	capture(t, func() error {
		_, err := load(gcfg, false)
		if err == nil {
			return fmt.Errorf("should fail")
		}

		return nil
	})
}

func capture(t *testing.T, fn func() error) string {
	t.Helper()

	old := os.Stdout

	oldcol := color.NoColor
	color.NoColor = true

	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	done := make(chan string)
	go func() {
		buf := &bytes.Buffer{}
		_, _ = io.Copy(buf, r)
		done <- buf.String()
	}()

	err = fn()
	// back to normal
	_ = w.Close()
	os.Stdout = old
	color.NoColor = oldcol
	assert.NoError(t, err)
	out := <-done

	return strings.TrimSpace(out)
}
