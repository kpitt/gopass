package age

import (
	"context"
	"fmt"

	"github.com/kpitt/gopass/internal/backend"
	"github.com/kpitt/gopass/pkg/debug"
)

const (
	name = "age"
)

func init() {
	backend.CryptoRegistry.Register(backend.Age, name, &loader{})
}

type loader struct{}

func (l loader) New(ctx context.Context) (backend.Crypto, error) {
	debug.Log("Using Crypto Backend: %s", name)

	return New()
}

func (l loader) Handles(ctx context.Context, s backend.Storage) error {
	if s.Exists(ctx, IDFile) {
		return nil
	}

	return fmt.Errorf("not supported")
}

func (l loader) Priority() int {
	return 10
}

func (l loader) String() string {
	return name
}
