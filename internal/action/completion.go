package action

import (
	"fmt"
	"regexp"
	"runtime"

	fishcomp "github.com/kpitt/gopass/internal/completion/fish"
	zshcomp "github.com/kpitt/gopass/internal/completion/zsh"
	"github.com/kpitt/gopass/internal/out"
	"github.com/kpitt/gopass/internal/tree"
	"github.com/kpitt/gopass/pkg/ctxutil"
	"github.com/urfave/cli/v2"
)

var escapeRegExp = regexp.MustCompile(`('|"|\s|\(|\)|\<|\>|\&|\;|\#|\\|\||\*|\?)`)

// bashEscape Escape special characters with `\`.
func bashEscape(s string) string {
	return escapeRegExp.ReplaceAllStringFunc(s, func(c string) string {
		if c == `\` {
			return `\\\\`
		}

		if c == `'` {
			return `\` + c
		}

		if c == `"` {
			return `\\\` + c
		}

		return `\\` + c
	})
}

// Complete prints a list of all password names to os.Stdout.
func (s *Action) Complete(c *cli.Context) {
	ctx := ctxutil.WithGlobalFlags(c)
	_, err := s.Store.IsInitialized(ctx) // important to make sure the structs are not nil.
	if err != nil {
		out.Errorf(ctx, "Store not initialized: %s", err)

		return
	}
	list, err := s.Store.List(ctx, tree.INF)
	if err != nil {
		return
	}

	for _, v := range list {
		fmt.Fprintln(stdout, bashEscape(v))
	}
}

// CompletionBash returns a bash script used for auto completion.
func (s *Action) CompletionBash(c *cli.Context) error {
	out := `_gopass_bash_autocomplete() {
     local cur opts base
     COMPREPLY=()
     cur="${COMP_WORDS[COMP_CWORD]}"
     opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
     local IFS=$'\n'
     COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
     return 0
 }

`
	out += "complete -F _gopass_bash_autocomplete " + s.Name
	if runtime.GOOS == "windows" {
		out += "\ncomplete -F _gopass_bash_autocomplete " + s.Name + ".exe"
	}
	fmt.Fprintln(stdout, out)

	return nil
}

// CompletionFish returns an autocompletion script for fish.
func (s *Action) CompletionFish(a *cli.App) error {
	if a == nil {
		return fmt.Errorf("app is nil")
	}
	comp, err := fishcomp.GetCompletion(a)
	if err != nil {
		return err
	}

	fmt.Fprintln(stdout, comp)

	return nil
}

// CompletionZSH returns a zsh completion script.
func (s *Action) CompletionZSH(a *cli.App) error {
	comp, err := zshcomp.GetCompletion(a)
	if err != nil {
		return err
	}

	fmt.Fprintln(stdout, comp)

	return nil
}
