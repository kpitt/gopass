# `gopass` command

Calling `gopass` without any command argument is a common entry point and
has two different modes.

## Synopsis

```
$ gopass
$ gopass entry
$ gopass -c entry
```

## Modes of operation

* Invoked without any arguments `gopass` will list all available secrets (equivalent to `list` command).
* Invoked with one argument it will perform a (fuzzy) search and display a list of matches or the secret directly (if exactly one match).
* Invoked with two arguments it will do search and if there is a match display the named key.

## Flags

Note: DO NOT use in scripts! Use `gopass show` instead.

Flag | Aliases | Description
---- | ------- | -----------
`--clip` | `-c` | Copy the password value into the clipboard and don't show the content.
`--yes` |  | Assume yes on all yes/no questions or use the default on all others.

