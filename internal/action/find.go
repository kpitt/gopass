package action

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/kpitt/gopass/internal/action/exit"
	"github.com/kpitt/gopass/internal/cui"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/internal/tree"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/kpitt/gopass/pkg/debug"
	"github.com/schollz/closestmatch"
	"github.com/urfave/cli/v2"
)

// Find runs the find command action, without fuzzy search.
func (s *Action) Find(c *cli.Context) error {
	if !c.Args().Present() {
		return exit.Error(exit.Usage, nil, "Usage: %s find <pattern>", s.Name)
	}

	ctx := ctxutil.WithGlobalFlags(c)

	return s.find(ctx, c, c.Args().First(), nil, false)
}

// FindFuzzy runs a fuzzy find of the specified name and attempts to show the result.
func (s *Action) FindFuzzy(c *cli.Context, name string) error {
	return s.find(c.Context, c, name, s.show, true)
}

// see action.show - context, cli context, name, key, rescurse.
type showFunc func(context.Context, *cli.Context, string, bool) error

func (s *Action) find(ctx context.Context, c *cli.Context, needle string, cb showFunc, fuzzy bool) error {
	// get all existing entries.
	haystack, err := s.Store.List(ctx, tree.INF)
	if err != nil {
		return exit.Error(exit.List, err, "failed to list store: %s", err)
	}

	// filter our the ones from the haystack matching the needle.
	needle = strings.ToLower(needle)
	choices := filter(haystack, needle)

	// if we have an exact match print it.
	if len(choices) == 1 {
		if cb == nil {
			out.Printf(ctx, choices[0])

			return nil
		}
		out.OKf(ctx, "Found exact match in %q", choices[0])

		return cb(ctx, c, choices[0], false)
	}

	// if we don't have a match yet try a fuzzy search.
	if len(choices) < 1 && fuzzy {
		// try fuzzy match.
		cm := closestmatch.New(haystack, []int{2})
		choices = cm.ClosestN(needle, 5)
	}

	// if there are still no results we abort.
	if len(choices) < 1 {
		return exit.Error(exit.NotFound, nil, "no results found")
	}

	// do not invoke wizard if not printing to terminal or if
	// gopass find/search was invoked directly (for scripts).
	if !ctxutil.IsTerminal(ctx) || (c != nil && c.Command.Name == "find") {
		for _, value := range choices {
			out.Printf(ctx, value)
		}

		return nil
	}

	return s.findSelection(ctx, c, choices, needle, cb)
}

// findSelection runs a wizard that lets the user select an entry.
func (s *Action) findSelection(ctx context.Context, c *cli.Context, choices []string, needle string, cb showFunc) error {
	if cb == nil {
		return fmt.Errorf("callback is nil")
	}
	if len(choices) < 1 {
		return fmt.Errorf("out of options")
	}

	sort.Strings(choices)
	act, sel := cui.GetSelection(ctx, "Found secrets - Please select an entry", choices)
	debug.Log("Action: %s - Selection: %d", act, sel)

	switch act {
	case "default":
		// display or copy selected entry.
		fmt.Fprintln(stdout, choices[sel])

		return cb(ctx, c, choices[sel], false)
	case "copy":
		// display selected entry.
		fmt.Fprintln(stdout, choices[sel])

		return cb(WithClip(ctx, true), c, choices[sel], false)
	case "show":
		// display selected entry.
		fmt.Fprintln(stdout, choices[sel])

		return cb(WithClip(ctx, false), c, choices[sel], false)
	case "sync":
		// run sync and re-run show/find workflow.
		if err := s.Sync(c); err != nil {
			return err
		}

		return cb(ctx, c, needle, true)
	case "edit":
		// edit selected entry.
		fmt.Fprintln(stdout, choices[sel])

		return s.edit(ctx, c, choices[sel])
	default:
		return exit.Error(exit.Aborted, nil, "user aborted")
	}
}

func filter(l []string, needle string) []string {
	choices := make([]string, 0, 10)
	for _, value := range l {
		if strings.Contains(strings.ToLower(value), needle) {
			choices = append(choices, value)
		}
	}

	return choices
}
