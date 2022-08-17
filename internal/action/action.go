package action

import (
	"io"
	"os"
	"path/filepath"

	"github.com/kpitt/gopass/internal/config"
	"github.com/kpitt/gopass/internal/reminder"
	"github.com/kpitt/gopass/internal/store/root"
	"github.com/kpitt/gopass/pkg/debug"
)

var (
	stdin  io.Reader = os.Stdin
	stdout io.Writer = os.Stdout
)

// Action knows everything to run gopass CLI actions.
type Action struct {
	Name    string
	Store   *root.Store
	cfg     *config.Config
	version string
	rem     *reminder.Store
}

// New returns a new Action wrapper.
func New(cfg *config.Config, version string) (*Action, error) {
	return newAction(cfg, version, true)
}

func newAction(cfg *config.Config, version string, remind bool) (*Action, error) {
	name := "gopass"
	if len(os.Args) > 0 {
		name = filepath.Base(os.Args[0])
	}

	act := &Action{
		Name:    name,
		cfg:     cfg,
		version: version,
		Store:   root.New(cfg),
	}

	if remind {
		r, err := reminder.New()
		if err != nil {
			debug.Log("failed to init reminder: %s", err)
		} else {
			// only populate the reminder variable on success, the implementation.
			// can handle being called on a nil pointer.
			act.rem = r
		}
	}

	return act, nil
}

// String implement fmt.Stringer.
func (s *Action) String() string {
	return s.Store.String()
}
