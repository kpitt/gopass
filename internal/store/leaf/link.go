package leaf

import (
	"context"
	"errors"
	"fmt"

	"github.com/kpitt/gopass/internal/queue"
	"github.com/kpitt/gopass/internal/store"
	"github.com/kpitt/gopass/pkg/debug"
)

// Link creates a symlink.
func (s *Store) Link(ctx context.Context, from, to string) error {
	if !s.Exists(ctx, from) {
		return fmt.Errorf("Source %q does not exists", from)
	}

	if s.Exists(ctx, to) {
		return fmt.Errorf("Destination %q already exists", to)
	}

	if err := s.storage.Link(ctx, s.Passfile(from), s.Passfile(to)); err != nil {
		return fmt.Errorf("Failed to create symlink from %q to %q: %w", from, to, err)
	}

	debug.Log("created symlink from %q to %q", from, to)

	if err := s.storage.Add(ctx, s.Passfile(to)); err != nil {
		if errors.Is(err, store.ErrGitNotInit) {
			return nil
		}

		return fmt.Errorf("Failed to add %q to git: %w", to, err)
	}

	// try to enqueue this task, if the queue is not available
	// it will return the task and we will execute it inline
	t := queue.GetQueue(ctx).Add(func(ctx context.Context) error {
		return s.gitCommitAndPush(ctx, to)
	})

	return t(ctx)
}
