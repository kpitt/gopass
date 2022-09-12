package action

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Version prints the gopass version.
func (s *Action) Version(c *cli.Context) error {
	cli.VersionPrinter(c)
	fmt.Fprintln(stdout)

	return nil
}
