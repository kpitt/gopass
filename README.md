# gopass

[![Build Status](https://img.shields.io/github/workflow/status/kpitt/gopass/Build%20gopass/master)](https://github.com/kpitt/gopass/actions/workflows/build.yml?query=branch%3Amaster)
[![Go Report Card](https://goreportcard.com/badge/github.com/kpitt/gopass)](https://goreportcard.com/report/github.com/kpitt/gopass)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/kpitt/gopass/blob/master/LICENSE)
[![Github All Releases](https://img.shields.io/github/downloads/kpitt/gopass/total.svg)](https://github.com/kpitt/gopass/releases)

This is an experimental fork of the [gopass](https://github.com/gopasspw/gopass) password manager. I really like the functionality provided by `gopass`, but the interface just doesn't _feel_ right for me personally. I think it's just a bit too "playful" for my tastes. As part of my journey to learn more about the Go language, I decided to experiment with `gopass` to see if I can create a cleaner, more professional interface along the lines of the [GitHub CLI Design Guidelines](https://primer.style/cli/).

I've also been getting more and more frustrated with how cumbersome GPG encryption is to use. I've become quite interested in the [age](https://age-encryption.org) encryption algorithm used in [passage](https://github.com/FiloSottile/passage), a fork of the original ZX2C4 [pass](https://www.passwordstore.org/) password manager. `gopass` already has some support for age encryption, but I think it would also be interesting to try to implement full interoperability with `passage`.

Note that updating this README will probably be one of the last things to get done, so the following information may be out-of-sync with the actual implementation for quite some time.

## Introduction

gopass is a password manager for the command line written in Go.
It works on all major desktop and server operating
systems (Linux, MacOS, BSD, Windows).

For detailed usage and installation instructions please check out our [documentation](docs/).

## Features

Please see [docs/features.md](https://github.com/kpitt/gopass/blob/master/docs/features.md) for an extensive list of all features along with several usage examples. Some examples are available in our
[example password store](https://github.com/gopasspw/password-store-example).

| **Feature**                 | **State**     | **Description**                                                   |
| --------------------------- | ------------- | ----------------------------------------------------------------- |
| Secure secret storage       | *stable*      | Securely storing encrypted secrets                                |
| Multiple stores             | *stable*      | Mount multiple stores in your root store, like file systems       |
| Recipient management        | *stable*      | Easily manage multiple users of each store                        |
| password quality assistance | *beta*        | Checks existing or new passwords for common flaws **offline**     |
| password leak checker       | *integration* | Perform **offline** checks against known leaked passwords using [gopass-hibp](https://github.com/gopasspw/gopass-hibp)  |
| PAGER support               | *stable*      | Automatically invoke a pager on long output                       |
| JSON API                    | *integration* | Allow gopass to be used as a native extension for browser plugins |
| Automatic fuzzy search      | *stable*      | Automatically search for matching store entries if a literal entry was not found |
| gopass sync                 | *stable*      | Easy to use syncing of remote repos and GPG keys                  |
| Desktop Notifications       | *stable*      | Display desktop notifications and completing long running operations |
| REPL                        | *beta*        | Integrated Read-Eval-Print-Loop shell with autocompletion by running `gopass`. |
| OTP support                 | *stable*      | Generate TOTP/(HOTP) tokens based on the stored secret            |
| Extensions                  |               | Extend gopass with custom commands using our API                  |
| Fully open source!          |               | No need to trust it, check our code and/or improve it!            |

## Design Principles

Gopass is a versatile command line based password manager that is being developed with the following principles in mind:

- **Easy**: For technical users (i.e. those who are used to the command line) it should be easy to get started with gopass.
- **Secure**: Security is hard. We aim to make it as easy as possible while still providing a good level of protection against common adversaries. *Caution*: If your personal threat level is very high, we might not offer a good tool for you.
- **Extensible**: While Gopass includes a fair amount of useful features, we can't cover every use-case. To support more special use cases we want to provide a clean and simple API to integration gopass into your own binaries.

## Installation

Please see [docs/setup.md](https://github.com/kpitt/gopass/blob/master/docs/setup.md).

If you have [Go](https://golang.org/) 1.18 (or greater) installed:

```bash
go install github.com/kpitt/gopass@latest
```
(and make sure your `$GOBIN` is in your `$PATH`.)

WARNING: Please prefer releases, unless you want to contribute to the
development of gopass. The master branch might not be stable and can contain breaking changes without any notice.

## Getting Started

Either initialize a new git repository or clone an existing one.

### New password store

```
$ gopass setup

🌟 Welcome to gopass!
- Initializing a new password store...
- Configuring your password store...
? Please select a private key for encrypting secrets:
[0] gpg - 0xFEEDBEEF - John Doe <john.doe@example.org>
Please enter the number of a key (0-12, [q]uit) (q to abort) [0]: 0
? Do you want to add a git remote? [y/N/q]: y
Configuring the git remote...
Please enter the git remote for your shared store []: git@gitlab.example.org:john/passwords.git
✓ Configured
```

Hint: `gopass setup` will use `gpg` encryption and `git` storage by default.

### Existing password store

```
$ gopass clone git@gitlab.example.org:john/passwords.git

🌟 Welcome to gopass!
- Cloning an existing password store from "git@gitlab.example.org:john/passwords.git"...
! Cloning git repository "git@gitlab.example.org:john/passwords.git" to "/home/john/.local/share/gopass/stores/root"...
! Configuring git repository...
- Gathering information for the git repository...
? What is your name? [John Doe]:
? What is your email? [john.doe@example.org]:
Your password store is ready to use! Have a look around: `gopass list`
```

## Upgrade

To use the self-updater run:
```bash
gopass update
```

or to upgrade with Go installed, run:
```bash
go install github.com/kpitt/gopass@latest
```

Otherwise, use the setup docs mentioned in the installation section to reinstall the latest version.

## Development

This project uses [GitHub Flow](https://guides.github.com/introduction/flow/). In other words, create feature branches from master, open an PR against master, and rebase onto master if necessary.

We aim for compatibility with the [latest stable Go Release](https://golang.org/dl/) only.

While this project is maintained by volunteers in their free time we aim to triage issues weekly and release a new version at least every quarter.

## Credit & License

gopass is licensed under the terms of the MIT license. You can find the complete text in `LICENSE`.

Please refer to the Git commit log for a complete list of contributors.

## Community

gopass is developed in the open. Here are some of the channels we use to communicate and contribute:

* Issue tracker: Use the [GitHub issue tracker](https://github.com/kpitt/gopass/issues) to file bugs and feature requests.

## Integrations

- [gopassbridge](https://github.com/gopasspw/gopassbridge): Browser plugin for Firefox, Chrome and other Chromium based browsers
- [gopass-ui](https://github.com/codecentric/gopass-ui): Graphical user interface for gopass
- [kubectl gopass](https://github.com/gopasspw/kubectl-gopass): Kubernetes / kubectl plugin to support reading and writing secrets directly from/to gopass.
- [gopass alfred](https://github.com/gopasspw/gopass-alfred): Alfred workflow to use gopass from the Alfred Mac launcher
- [git-credential-gopass](https://github.com/gopasspw/git-credential-gopass): Integrate gopass as an git-credential helper
- [gopass-hibp](https://github.com/gopasspw/gopass-hibp): haveibeenpwned.com leak checker
- [gopass-jsonapi](https://github.com/gopasspw/gopass-jsonapi): native messaging for browser plugins, e.g. gopassbridge
- [gopass-summon-prover](https://github.com/gopasspw/gopass-summon-provider): gopass as a summon provider
- [`terraform-provider-gopass`](https://github.com/camptocamp/terraform-provider-pass): a Terraform provider to interact with gopass
- [chezmoi](https://github.com/twpayne/chezmoi): dotfile manager with gopass support
- [tessen](https://github.com/ayushnix/tessen): autotype and copy gopass data on wayland compositors on Linux
- [raycast-gopass](https://github.com/raycast/extensions/tree/main/extensions/gopass): a gopass extension for Raycast Mac launcher

## Mobile apps

- [Pass - Password Store](https://apps.apple.com/us/app/pass-password-store/id1205820573) - iOS, [source code](https://github.com/mssun/passforios), [supports only 1 repository now](https://github.com/mssun/passforios/issues/88)
- [Password Store](https://play.google.com/store/apps/details?id=dev.msfjarvis.aps) - Android

## Related Projects

- [pass](https://www.passwordstore.org) - The inspiration for this project, by Jason A. Donenfeld. `gopass` is a drop-in replacement for `pass` and can be used interchangeably (mostly!).
- [passage](https://github.com/FiloSottile/passage) - passage is a fork of [password-store](https://www.passwordstore.org) that uses
[age](https://age-encryption.org) as a backend instead of GnuPG. `gopass` has some amount of support for `passage` but can not be used fully interchangeably as of today. This might change in the future.

## Contributing

We welcome any contributions. Please see the [CONTRIBUTING.md](https://github.com/gopasspw/gopass/blob/master/CONTRIBUTING.md) file for instructions on how to submit changes.

## Further Documentation

* [Security, Known Limitations, and Caveats](https://github.com/gopasspw/gopass/blob/master/docs/security.md)
* [Configuration](https://github.com/gopasspw/gopass/blob/master/docs/config.md)
* [FAQ](https://github.com/gopasspw/gopass/blob/master/docs/faq.md)
* [JSON API](https://github.com/gopasspw/gopass-jsonapi)
* [Gopass as Summon provider](https://github.com/gopasspw/gopass-summon-provider)

## External Documentation
* [gopass cheat sheet](https://woile.github.io/gopass-cheat-sheet/) ([source on github](https://github.com/Woile/gopass-cheat-sheet))
* [gopass presentation](https://woile.github.io/gopass-presentation/) ([source on github](https://github.com/Woile/gopass-presentation))
