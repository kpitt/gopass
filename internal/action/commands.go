package action

import (
	"fmt"

	"github.com/kpitt/gopass/internal/backend"
	"github.com/kpitt/gopass/pkg/debug"
	"github.com/urfave/cli/v2"
)

// ShowFlags returns the flags for the show command. Exported to re-use in main
// for the default command.
func ShowFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "yes",
			Aliases: []string{"y"},
			Usage:   "Always answer yes to yes/no questions",
		},
		&cli.BoolFlag{
			Name:    "clip",
			Aliases: []string{"c"},
			Usage:   "Copy the password value into the clipboard",
		},
		&cli.BoolFlag{
			Name:  "qr",
			Usage: "Print the password as a QR Code",
		},
		&cli.BoolFlag{
			Name:    "password",
			Aliases: []string{"o"},
			Usage:   "Display only the password. Takes precedence over all other flags.",
		},
		&cli.StringFlag{
			Name:    "revision",
			Aliases: []string{"r"},
			Usage:   "Show a past revision. Does NOT support Git shortcuts. Use exact revision or -<N> to select the Nth oldest revision of this entry.",
		},
		&cli.BoolFlag{
			Name:    "noparsing",
			Aliases: []string{"n"},
			Usage:   "Do not parse the output.",
		},
		&cli.StringFlag{
			Name:  "chars",
			Usage: "Print specific characters from the secret",
		},
	}
}

// GetCommands returns the cli commands exported by this module.
func (s *Action) GetCommands() []*cli.Command {
	cmds := []*cli.Command{
		{
			Name:      "audit",
			Usage:     "Decrypt all secrets and scan for weak or leaked passwords",
			ArgsUsage: "[filter]",
			Description: "" +
				"This command decrypts all secrets and checks for common flaws and (optionally) " +
				"against a list of previously leaked passwords.",
			Before: s.IsInitialized,
			Action: s.Audit,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "expiry",
					Usage: "Age in days before a password is considered expired. Setting this will only check expiration.",
				},
			},
		},
		{
			Name:      "cat",
			Usage:     "Decode and print content of a binary secret to stdout, or encode and insert from stdin",
			ArgsUsage: "[secret]",
			Description: "" +
				"This command is similar to the way cat works on the command line. " +
				"It can either be used to retrieve the decoded content of a secret " +
				"similar to 'cat file' or vice versa to encode the content from STDIN " +
				"to a secret.",
			Before:       s.IsInitialized,
			Action:       s.Cat,
			BashComplete: s.Complete,
		},
		{
			Name:      "clone",
			Usage:     "Clone a password store from a git repository",
			ArgsUsage: "[git-repo] [mount-point]",
			Description: "" +
				"This command clones an existing password store from a git remote to " +
				"a local password store. Can be either used to initialize a new root store " +
				"or to add a new mounted sub-store. " +
				"" +
				"Needs at least one argument (git URL) to clone from. " +
				"Accepts a second argument (mount location) to clone and mount a sub-store, e.g. " +
				"'gopass clone git@example.com/store.git foo/bar'",
			Action: s.Clone,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "path",
					Usage: "Path to clone the repo to",
				},
				&cli.StringFlag{
					Name:  "crypto",
					Usage: fmt.Sprintf("Select crypto backend %v", backend.CryptoRegistry.BackendNames()),
				},
				&cli.BoolFlag{
					Name:  "check-keys",
					Usage: "Check for valid decryption keys. Generate new keys if none are found.",
					Value: true,
				},
			},
		},
		{
			Name:      "config",
			Usage:     "Display and edit the configuration file",
			ArgsUsage: "[key [value]]",
			Description: "" +
				"This command allows for easy printing and editing of the configuration. " +
				"Without argument, the entire config is printed. " +
				"With a single argument, a single key can be printed. " +
				"With two arguments a setting specified by key can be set to value.",
			Action:       s.Config,
			BashComplete: s.ConfigComplete,
		},
		{
			Name:      "copy",
			Aliases:   []string{"cp"},
			Usage:     "Copy secrets from one location to another",
			ArgsUsage: "[from] [to]",
			Description: "" +
				"This command copies an existing secret in the store to another location. " +
				"This also works across different sub-stores. If the source is a directory it will " +
				"automatically copy recursively. In that case, the source directory is re-created " +
				"at the destination if no trailing slash is found, otherwise the contents are " +
				"flattened (similar to rsync).",
			Before:       s.IsInitialized,
			Action:       s.Copy,
			BashComplete: s.Complete,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "Force to copy the secret and overwrite existing one",
				},
			},
		},
		{
			Name:      "create",
			Aliases:   []string{"new"},
			Usage:     "Easy creation of new secrets",
			ArgsUsage: "[secret]",
			Description: "" +
				"This command starts a wizard to aid in creation of new secrets.",
			Before: s.IsInitialized,
			Action: s.Create,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "store",
					Aliases: []string{"s"},
					Usage:   "Which store to use",
				},
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "Force path selection",
				},
			},
		},
		{
			Name:      "delete",
			Usage:     "Remove one or many secrets from the store",
			ArgsUsage: "[secret [key]]",
			Description: "" +
				"This command removes secrets. It can work recursively on folders. " +
				"Recursing across stores is purposefully not supported.",
			Aliases:      []string{"remove", "rm"},
			Before:       s.IsInitialized,
			Action:       s.Delete,
			BashComplete: s.Complete,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "recursive",
					Aliases: []string{"r"},
					Usage:   "Recursive delete files and folders",
				},
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "Force to delete the secret",
				},
			},
		},
		{
			Name:      "edit",
			Usage:     "Edit new or existing secrets",
			ArgsUsage: "[secret]",
			Description: "" +
				"Use this command to insert a new secret or edit an existing one using " +
				"your $EDITOR. It will attempt to create a secure temporary directory " +
				"for storing your secret while the editor is accessing it. Please make " +
				"sure your editor doesn't leak sensitive data to other locations while " +
				"editing.\n" +
				"Note: If $EDITOR is not set we will try 'editor'. If that's not available " +
				"either we fall back to 'vi'. Consider using 'update-alternatives --config editor " +
				"to change the defaults.",
			Before:       s.IsInitialized,
			Action:       s.Edit,
			Aliases:      []string{"set"},
			BashComplete: s.Complete,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "editor",
					Aliases: []string{"e"},
					Usage:   "Use this editor binary",
				},
				&cli.BoolFlag{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "Create a new secret if none found",
				},
			},
		},
		{
			Name:      "find",
			Usage:     "Search for secrets",
			ArgsUsage: "<pattern>",
			Description: "" +
				"List all secrets that match the specified search pattern.",
			Before:       s.IsInitialized,
			Action:       s.Find,
			Aliases:      []string{"search"},
			BashComplete: s.Complete,
		},
		{
			Name:      "fsck",
			Usage:     "Check store integrity",
			ArgsUsage: "[filter]",
			Description: "" +
				"Check the integrity of the given sub-store or all stores if none are specified. " +
				"Will automatically fix all issues found.",
			Before:       s.IsInitialized,
			Action:       s.Fsck,
			BashComplete: s.MountsComplete,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "decrypt",
					Usage: "Decrypt and reencryt during fsck.",
				},
			},
		},
		{
			Name:      "fscopy",
			Usage:     "Copy files from or to the password store",
			ArgsUsage: "[from] [to]",
			Description: "" +
				"This command either reads a file from the filesystem and writes the " +
				"encoded and encrypted version in the store or it decrypts and decodes " +
				"a secret and writes the result to a file. Either source or destination " +
				"must be a file and the other one a secret. If you want the source to " +
				"be securely removed after copying, use 'gopass binary move'",
			Before:       s.IsInitialized,
			Action:       s.BinaryCopy,
			BashComplete: s.Complete,
			Hidden:       true,
		},
		{
			Name:      "fsmove",
			Usage:     "Move files from or to the password store",
			ArgsUsage: "[from] [to]",
			Description: "" +
				"This command either reads a file from the filesystem and writes the " +
				"encoded and encrypted version in the store or it decrypts and decodes " +
				"a secret and writes the result to a file. Either source or destination " +
				"must be a file and the other one a secret. The source will be wiped " +
				"from disk or from the store after it has been copied successfully " +
				"and validated. If you don't want the source to be removed use " +
				"'gopass binary copy'",
			Before:       s.IsInitialized,
			Action:       s.BinaryMove,
			BashComplete: s.Complete,
			Hidden:       true,
		},
		{
			Name:      "generate",
			Usage:     "Generate a new password",
			ArgsUsage: "[secret [key [length]|length]]",
			Description: "" +
				"Dialog to generate a new password and write it into a new or existing secret. " +
				"By default, the new password will replace the first line of an existing secret (or create a new one).",
			Before:       s.IsInitialized,
			Action:       s.Generate,
			BashComplete: s.CompleteGenerate,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "clip",
					Aliases: []string{"c"},
					Usage:   "Copy the generated password to the clipboard",
				},
				&cli.BoolFlag{
					Name:    "print",
					Aliases: []string{"p"},
					Usage:   "Print the generated password to the terminal",
				},
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "Force to overwrite existing password",
				},
				&cli.BoolFlag{
					Name:    "edit",
					Aliases: []string{"e"},
					Usage:   "Open secret for editing after generating a password",
				},
				&cli.BoolFlag{
					Name:    "symbols",
					Aliases: []string{"s"},
					Usage:   "Use symbols in the password",
				},
				&cli.StringFlag{
					Name:    "generator",
					Aliases: []string{"g"},
					Usage:   "Choose a password generator, use one of: cryptic, memorable, xkcd or external. Default: cryptic",
				},
				&cli.BoolFlag{
					Name:  "strict",
					Usage: "Require strict character class rules",
				},
				&cli.StringFlag{
					Name:    "sep",
					Aliases: []string{"xkcdsep", "xs"},
					Usage:   "Word separator for generated passwords. If no separator is specified, the words are combined without spaces/separator and the first character of words is capitalised.",
					Value:   "",
				},
				&cli.StringFlag{
					Name:    "lang",
					Aliases: []string{"xkcdlang", "xl"},
					Usage:   "Language to generate password from, currently only en (english, default) is supported",
					Value:   "en",
				},
			},
		},
		{
			Name:  "git",
			Usage: "Run a git command inside a password store",
			Description: `If the password store is a git repository, execute a git command in the password store directory.

Use the "git init" command if the store does not yet have a git repository.`,
			ArgsUsage: "[git-command-args...]",
			Before:    s.IsInitialized,
			Action:    s.Git,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "store",
					Usage: "Store to operate on",
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:        "init",
					Usage:       "Initialize git repository",
					Description: "Initialize and configure a git repository inside an existing password store.",
					Before:      s.IsInitialized,
					Action:      s.RCSInit,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "store",
							Usage: "Store to operate on",
						},
						&cli.StringFlag{
							Name:  "name",
							Usage: "Git user name",
						},
						&cli.StringFlag{
							Name:  "email",
							Usage: "Git user email",
						},
					},
				},
			},
		},
		{
			Name:      "grep",
			Usage:     "Search for secrets files containing search-string when decrypted.",
			ArgsUsage: "[needle]",
			Description: "" +
				"This command decrypts all secrets and performs a pattern matching on the " +
				"content.",
			Before: s.IsInitialized,
			Action: s.Grep,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "regexp",
					Aliases: []string{"r"},
					Usage:   "Interpret pattern as RE2 regular expression",
				},
			},
		},
		{
			Name:      "history",
			Usage:     "Show password history",
			ArgsUsage: "[secret]",
			Aliases:   []string{"hist"},
			Description: "" +
				"Display the change history for a secret",
			Before:       s.IsInitialized,
			Action:       s.History,
			BashComplete: s.Complete,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Include passwords in output",
				},
			},
		},
		{
			Name:      "init",
			Usage:     "Initialize new password store.",
			ArgsUsage: "[gpg-id]",
			Description: "" +
				"Initialize new password storage and use gpg-id for encryption.",
			Action: s.Init,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "path",
					Aliases: []string{"p"},
					Usage:   "Set the sub-store path to operate on",
				},
				&cli.StringFlag{
					Name:    "store",
					Aliases: []string{"s"},
					Usage:   "Set the name of the sub-store",
				},
				&cli.StringFlag{
					Name:  "crypto",
					Usage: fmt.Sprintf("Select crypto backend %v", backend.CryptoRegistry.BackendNames()),
					Value: "gpgcli",
				},
				&cli.StringFlag{
					Name:  "storage",
					Usage: fmt.Sprintf("Select storage backend %v", backend.StorageRegistry.BackendNames()),
					Value: "gitfs",
				},
				&cli.StringFlag{
					Name:    "remote",
					Aliases: []string{"R"},
					Usage:   "URL of remote Git repository for this store",
				},
			},
		},
		{
			Name:      "insert",
			Usage:     "Insert a new secret",
			ArgsUsage: "[secret]",
			Description: "" +
				"Insert a new secret. Optionally, echo the secret back to the console during entry. " +
				"Or, optionally, the entry may be multiline. " +
				"Prompt before overwriting existing secret unless forced.",
			Before:       s.IsInitialized,
			Action:       s.Insert,
			BashComplete: s.Complete,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "echo",
					Aliases: []string{"e"},
					Usage:   "Display secret while typing",
				},
				&cli.BoolFlag{
					Name:    "multiline",
					Aliases: []string{"m"},
					Usage:   "Insert using $EDITOR",
				},
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "Overwrite any existing secret and do not prompt to confirm recipients",
				},
				&cli.BoolFlag{
					Name:    "append",
					Aliases: []string{"a"},
					Usage:   "Append data read from STDIN to existing data",
				},
			},
		},
		{
			Name:      "link",
			Usage:     "Create a symlink",
			ArgsUsage: "[from] [to]",
			Description: "" +
				"This command creates a symlink from one entry in a mounted store to another entry. " +
				"Important: Does not cross mounts!",
			Aliases:      []string{"ln", "symlink"},
			Hidden:       true,
			Before:       s.IsInitialized,
			Action:       s.Link,
			BashComplete: s.Complete,
		},
		{
			Name:      "list",
			Usage:     "List existing secrets",
			ArgsUsage: "[prefix]",
			Description: "" +
				"This command will list all existing secrets. Provide a folder prefix to list " +
				"only certain subfolders of the store.",
			Aliases:      []string{"ls"},
			Before:       s.IsInitialized,
			Action:       s.List,
			BashComplete: s.Complete,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "limit",
					Aliases: []string{"l"},
					Usage:   "Display no more than this many levels of the tree",
				},
				&cli.BoolFlag{
					Name:    "flat",
					Aliases: []string{"f"},
					Usage:   "Print a flat list",
				},
				&cli.BoolFlag{
					Name:    "folders",
					Aliases: []string{"d"},
					Usage:   "Print a flat list of folders",
				},
				&cli.BoolFlag{
					Name:    "strip-prefix",
					Aliases: []string{"s"},
					Usage:   "Strip this prefix from filtered entries",
				},
			},
		},
		{
			Name:      "merge",
			Usage:     "Merge multiple secrets into one",
			ArgsUsage: "[to] [from]...",
			Description: "" +
				"This command implements a merge workflow to help deduplicate " +
				"secrets. It requires exactly one destination (may already exist) " +
				"and at least one source (must exist, can be multiple). gopass will " +
				"then merge all entries into one, drop into an editor, save the result " +
				"and remove all merged entries.",
			Before:       s.IsInitialized,
			Action:       s.Merge,
			BashComplete: s.Complete,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "Remove merged entries",
					Value:   true,
				},
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "Skip editor, merge entries unattended",
				},
			},
		},
		{
			Name:      "move",
			Aliases:   []string{"mv"},
			Usage:     "Move secrets from one location to another",
			ArgsUsage: "[from] [to]",
			Description: "" +
				"This command moves a secret from one path to another. This also works " +
				"across different sub-stores. If the source is a directory, the source directory " +
				"is re-created at the destination if no trailing slash is found, otherwise the " +
				"contents are flattened (similar to rsync).",
			Before:       s.IsInitialized,
			Action:       s.Move,
			BashComplete: s.Complete,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "Force to move the secret and overwrite existing one",
				},
			},
		},
		{
			Name:  "mounts",
			Usage: "Edit mounted stores",
			Description: "" +
				"This command displays all mounted password stores. It offers several " +
				"subcommands to create or remove mounts.",
			Before: s.IsInitialized,
			Action: s.MountsPrint,
			Subcommands: []*cli.Command{
				{
					Name:    "add",
					Aliases: []string{"mount"},
					Usage:   "Mount a password store",
					Description: "" +
						"This command allows for mounting an existing or new password store " +
						"at any path in an existing root store.",
					Before: s.IsInitialized,
					Action: s.MountAdd,
				},
				{
					Name:    "remove",
					Aliases: []string{"rm", "unmount", "umount"},
					Usage:   "Umount an mounted password store",
					Description: "" +
						"This command allows to unmount an mounted password store. This will " +
						"only updated the configuration and not delete the password store.",
					Before:       s.IsInitialized,
					Action:       s.MountRemove,
					BashComplete: s.MountsComplete,
				},
				{
					Name:    "versions",
					Aliases: []string{"version"},
					Usage:   "Display mount provider versions",
					Description: "" +
						"This command displays version information of important external " +
						"commands used by the configured password store mounts.",
					Before: s.IsInitialized,
					Action: s.MountsVersions,
				},
			},
		},
		{
			Name:      "otp",
			Usage:     "Generate time- or hmac-based tokens",
			ArgsUsage: "[secret]",
			Aliases:   []string{"totp", "hotp"},
			Description: "" +
				"Tries to parse an OTP URL (otpauth://). URL can be TOTP or HOTP. " +
				"The URL can be provided on its own line or on a key value line with a key named 'totp'.",
			Before:       s.IsInitialized,
			Action:       s.OTP,
			BashComplete: s.Complete,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "clip",
					Aliases: []string{"c"},
					Usage:   "Copy the time-based token into the clipboard",
				},
				&cli.StringFlag{
					Name:    "qr",
					Aliases: []string{"q"},
					Usage:   "Write QR code to `FILE`",
				},
				&cli.BoolFlag{
					Name:    "continuous",
					Aliases: []string{"C"},
					Usage:   "Display tokens continuously until interrupted",
				},
			},
		},
		{
			Name:  "process",
			Usage: "Process a template file",
			Description: "" +
				"This command processes a template file. It will read the template file " +
				"and replace all variables with their values.",
			Before: s.IsInitialized,
			Action: s.Process,
		},
		{
			Name:  "recipients",
			Usage: "Edit recipient permissions",
			Description: "" +
				"This command displays all existing recipients for all mounted stores. " +
				"The subcommands allow adding or removing recipients.",
			Before: s.IsInitialized,
			Action: s.RecipientsPrint,
			Subcommands: []*cli.Command{
				{
					Name:    "add",
					Aliases: []string{"authorize"},
					Usage:   "Add any number of Recipients to any store",
					Description: "" +
						"This command adds any number of recipients to any existing store. " +
						"If none are given it will display a list of usable public keys. " +
						"After adding the recipient to the list it will re-encrypt the whole " +
						"affected store to make sure the recipient has access to all existing " +
						"secrets.",
					Before: s.IsInitialized,
					Action: s.RecipientsAdd,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "store",
							Usage: "Store to operate on",
						},
						&cli.BoolFlag{
							Name:  "force",
							Usage: "Force adding non-existing keys",
						},
					},
				},
				{
					Name:    "remove",
					Aliases: []string{"rm", "deauthorize"},
					Usage:   "Remove any number of Recipients from any store",
					Description: "" +
						"This command removes any number of recipients from any existing store. " +
						"If no recipients are provided, it will show a list of existing recipients " +
						"to choose from. It will refuse to remove the current user's key from the " +
						"store to avoid losing access. After removing the keys it will re-encrypt " +
						"all existing secrets. Please note that the removed recipients will still " +
						"be able to decrypt old revisions of the password store and any local " +
						"copies they might have. The only way to reliably remove a recipient is to " +
						"rotate all existing secrets.",
					Before:       s.IsInitialized,
					Action:       s.RecipientsRemove,
					BashComplete: s.RecipientsComplete,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "store",
							Usage: "Store to operate on",
						},
						&cli.BoolFlag{
							Name:  "force",
							Usage: "Force adding non-existing keys",
						},
					},
				},
			},
		},
		{
			Name:      "show",
			Usage:     "Display the content of a secret",
			ArgsUsage: "[secret]",
			Description: "" +
				"Show an existing secret and optionally put its first line on the clipboard. " +
				"If put on the clipboard, it will be cleared after 45 seconds.",
			Before:       s.IsInitialized,
			Action:       s.Show,
			BashComplete: s.Complete,
			Flags:        ShowFlags(),
		},
		{
			Name:      "sum",
			Usage:     "Compute the SHA256 checksum",
			ArgsUsage: "[secret]",
			Description: "" +
				"This command decodes an Base64 encoded secret and computes the SHA256 checksum " +
				"over the decoded data. This is useful to verify the integrity of an " +
				"inserted secret.",
			Aliases:      []string{"sha", "sha256"},
			Before:       s.IsInitialized,
			Action:       s.Sum,
			BashComplete: s.Complete,
		},
		{
			Name:  "sync",
			Usage: "Sync all local stores with their remotes",
			Description: "" +
				"Sync all local stores with their git remotes, if any, and check " +
				"any possibly affected gpg keys.",
			Before: s.IsInitialized,
			Action: s.Sync,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "store",
					Aliases: []string{"s"},
					Usage:   "Select the store to sync",
				},
			},
		},
		{
			Name:  "templates",
			Usage: "Edit templates",
			Description: "" +
				"List existing templates in the password store and allow for editing " +
				"and creating them.",
			Before: s.IsInitialized,
			Action: s.TemplatesPrint,
			Subcommands: []*cli.Command{
				{
					Name:         "show",
					Usage:        "Show a secret template.",
					Description:  "Display an existing template",
					Aliases:      []string{"cat"},
					Before:       s.IsInitialized,
					Action:       s.TemplatePrint,
					BashComplete: s.TemplatesComplete,
				},
				{
					Name:         "edit",
					Usage:        "Edit secret templates.",
					Description:  "Edit an existing or new template",
					Aliases:      []string{"create", "new"},
					Before:       s.IsInitialized,
					Action:       s.TemplateEdit,
					BashComplete: s.TemplatesComplete,
				},
				{
					Name:         "remove",
					Aliases:      []string{"rm"},
					Usage:        "Remove secret templates.",
					Description:  "Remove an existing template",
					Before:       s.IsInitialized,
					Action:       s.TemplateRemove,
					BashComplete: s.TemplatesComplete,
				},
			},
		},
		{
			Name:        "unclip",
			Usage:       "Internal command to clear clipboard",
			Description: "Clear the clipboard if the content matches the checksum.",
			Action:      s.Unclip,
			Hidden:      true,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "timeout",
					Usage: "Time to wait",
				},
				&cli.BoolFlag{
					Name:  "force",
					Usage: "Clear clipboard even if checksum mismatches",
				},
			},
		},
		{
			Name:  "version",
			Usage: "Display version",
			Description: "" +
				"This command displays version and build time information.",
			Action: s.Version,
		},
	}

	// crypto and storage backends can add their own commands if they need to
	for _, be := range backend.CryptoRegistry.Backends() {
		bc, ok := be.(commander)
		if !ok {
			debug.Log("Backend %s does not implement commander interface\n", be)

			continue
		}
		nc := bc.Commands()
		debug.Log("Backend %s added %d commands", be, len(nc))
		cmds = append(cmds, nc...)
	}

	return cmds
}

type commander interface {
	Commands() []*cli.Command
}
