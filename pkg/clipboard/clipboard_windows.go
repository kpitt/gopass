//go:build windows
// +build windows

package clipboard

import (
	"context"
	"os"
	"os/exec"
	"strconv"

	"github.com/kpitt/gopass/internal/pwschemes/argon2id"
	"github.com/kpitt/gopass/pkg/ctxutil"
)

// clear will spwan a copy of gopass that waits in a detached background
// process group until the timeout is expired. It will then compare the contents
// of the clipboard and erase it if it still contains the data gopass copied
// to it.
func clear(ctx context.Context, name string, content []byte, timeout int) error {
	hash, err := argon2id.Generate(string(content), 0)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, os.Args[0], "unclip", "--timeout", strconv.Itoa(timeout))
	cmd.Env = append(os.Environ(), "GOPASS_UNCLIP_NAME="+name)
	cmd.Env = append(cmd.Env, "GOPASS_UNCLIP_CHECKSUM="+hash)
	return cmd.Start()
}

func walkFn(int, func(int)) {}
