# Configuration

## Environment Variables

Some configuration options are only available through setting environment variables.

| **Option**                   | **Type** | **Description**                                                                                                  |
| ---------------------------- | -------- | ---------------------------------------------------------------------------------------------------------------- |
| `GOPASS_DEBUG`               | `bool`   | Set to any non-empty value to enable verbose debug output                                                        |
| `GOPASS_DEBUG_LOG`           | `string` | Set to a filename to enable debug logging                                                                        |
| `GOPASS_DEBUG_LOG_SECRETS`   | `bool`   | Set to any non-empty value to enable logging of credentials                                                      |
| `GOPASS_DEBUG_FUNCS`         | `string` | Comma separated filter for console debug output (functions)                                                      |
| `GOPASS_DEBUG_FILES`         | `string` | Comma separated filter for console debug output (files)                                                          |
| `GOPASS_UMASK`               | `octal`  | Set to any valid umask to mask bits of files created by gopass                                                   |
| `GOPASS_GPG_OPTS`            | `string` | Add any extra arguments, e.g. `--armor` you want to pass to GPG on every invocation                              |
| `GOPASS_EXTERNAL_PWGEN`      | `string` | Use an external password generator. See [Features](features.md#using-custom-password-generators) for details     |
| `GOPASS_CHARACTER_SET`       | `bool`   | Set to any non-empty value to restrict the characters used in generated passwords                                |
| `GOPASS_CONFIG`              | `string` | Set this to the absolute path to the configuration file                                                          |
| `GOPASS_HOMEDIR`             | `string` | Set this to the absolute path of the directory containing the `.config/` tree                                    |
| `GOPASS_NO_REMINDER`         | `bool`   | Set to any non-empty value to prevent reminders                                                                  |
| `GOPASS_CLIPBOARD_COPY_CMD`  | `string` | Use an external command to copy a password to the clipboard. See [GPaste](usecases/gpaste.md) for an example     |
| `GOPASS_CLIPBOARD_CLEAR_CMD` | `string` | Use an external command to remove a password from the clipboard. See [GPaste](usecases/gpaste.md) for an example |
| `GOPASS_GPG_BINARY` | `string` | Set this to the absolute path to the GPG binary if you need to override the value returned by `gpgconf`, e.g. [QubesOS](https://www.qubes-os.org/doc/split-gpg/). |
| `GOPASS_PW_DEFAULT_LENGTH`   | `int`    | Set to any integer value larger than zero to define a different default length in the `generate` command. By default the length is 24 characters. |
| `GOPASS_AUTOSYNC_INTERVAL` | `int` | Set this to the number of days between autosync runs. |
| `GOPASS_NO_AUTOSYNC` | `bool` | Set this to `true` to disable autosync. |

Variables not exclusively used by gopass

| **Option**             | **Type** | **Description**                                                                                        |
| ---------------------- | -------- | ------------------------------------------------------------------------------------------------------ |
| `PASSWORD_STORE_DIR`   | `string` | absolute path containing the password store (a directory). Only supported during initialization!       |
| `PASSWORD_STORE_UMASK` | `string` | Set to any valid umask to mask bits of files created by gopass (GOPASS_UMASK has precedence over this) |
| `EDITOR`               | `string` | command name to execute for editing password entries                                                   |
| `PAGER`                | `string` | the pager program used for `gopass list`. See [Features](features.md#auto-pager) for details           |
| `GIT_AUTHOR_NAME`      | `string` | name of the author, used by the git backend to create a commit                                         |
| `GIT_AUTHOR_EMAIL`     | `string` | email of the author, used by the git backend to create a commit                                        |
| `NO_COLOR`             | `bool`   | disable color output. See [no-color.org](https://no-color.org) for more information.                   |

## Configuration Options

During start up, gopass will look for a configuration file at `$HOME/.config/gopass/config.yml`. If one is not present, it will create one. If the config file already exists, it will attempt to parse it and load the settings. If this fails, the program will abort. Thus, if gopass is giving you trouble with a broken or incompatible configuration file, simply rename it or delete it.

All configuration options are also available for reading and writing through the sub-command `gopass config`.

* To display all values: `gopass config`
* To display a single value: `gopass config autoimport`
* To update a single value: `gopass config autoimport false`
* As many other sub-commands this command accepts a `--store` flag to operate on a given sub-store, provided the sub-store is a remote one.

This is a list of available options:

| **Option**       | **Type** | Description                                                                                                                                                                                    |
| ---------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `autoimport`     | `bool`   | Import missing keys stored in the pass repository without asking.                                                                                                                              |
| `cliptimeout`    | `int`    | How many seconds the secret is stored when using `-c`.                                                                                                                                         |
| `exportkeys`     | `bool`   | Export public keys of all recipients to the store.                                                                                                                                             |
| `nopager`        | `bool`   | Do not invoke a pager to display long lists.                                                                                                                                                   |
| `parsing`        | `bool`   | Enable parsing of output to have key-value and yaml secrets.                                                                                                                                   |
| `path`           | `string` | Path to the root store.                                                                                                                                                                        |
