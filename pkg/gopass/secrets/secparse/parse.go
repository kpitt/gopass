package secparse

import (
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/pkg/debug"
	"github.com/kpitt/gopass/pkg/gopass"
	"github.com/kpitt/gopass/pkg/gopass/secrets"
)

// Parse tries to parse a secret. It will start with the most specific
// secrets type.
//
//nolint:ireturn
func Parse(in []byte) (gopass.Secret, error) {
	var s gopass.Secret

	var err error

	s, err = secrets.ParseYAML(in)
	if err == nil {
		debug.Log("parsed as YAML: %+v", s)

		return s, nil
	}

	debug.Log("failed to parse as YAML: %s\n%s", err, out.Secret(string(in)))

	s, err = secrets.ParseKV(in)
	if err == nil {
		debug.Log("parsed as KV: %+v", s)

		return s, nil
	}

	debug.Log("failed to parse as KV: %s", err)

	s = secrets.ParsePlain(in)
	debug.Log("parsed as plain: %s", out.Secret(s.Bytes()))

	return s, nil
}
