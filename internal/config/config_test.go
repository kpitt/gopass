package config_test

import (
	"testing"

	_ "github.com/kpitt/gopass/internal/backend/crypto"
	_ "github.com/kpitt/gopass/internal/backend/storage"
	"github.com/kpitt/gopass/internal/config"
	"github.com/kpitt/gopass/tests/gptest"
	"github.com/stretchr/testify/assert"
)

func TestHomedir(t *testing.T) { //nolint:paralleltest
	assert.NotEqual(t, config.Homedir(), "")
}

func TestNewConfig(t *testing.T) { //nolint:paralleltest
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	cfg := config.New()
	cs := cfg.String()
	assert.Contains(t, cs, `&config.Config{AutoClip:false, AutoImport:false, ClipTimeout:45, ExportKeys:true, NoPager:false,`)
	assert.Contains(t, cs, `SafeContent:false, Mounts:map[string]string{},`)

	cfg = &config.Config{
		Mounts: map[string]string{
			"foo": "",
			"bar": "",
		},
	}
	cs = cfg.String()
	assert.Contains(t, cs, `&config.Config{AutoClip:false, AutoImport:false, ClipTimeout:0, ExportKeys:false, NoPager:false,`)
	assert.Contains(t, cs, `SafeContent:false, Mounts:map[string]string{"bar":"", "foo":""},`)
}

func TestSetConfigValue(t *testing.T) { //nolint:paralleltest
	u := gptest.NewUnitTester(t)
	defer u.Remove()

	cfg := config.New()
	assert.NoError(t, cfg.SetConfigValue("autoclip", "true"))
	assert.NoError(t, cfg.SetConfigValue("cliptimeout", "900"))
	assert.NoError(t, cfg.SetConfigValue("path", "/tmp"))
	assert.Error(t, cfg.SetConfigValue("autoclip", "yo"))
}
