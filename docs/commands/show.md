# `show` command

The `show` command is the most important and most frequently used command.
It allows displaying and copying the content of the secrets managed by gopass.

## Synopsis

```
$ gopass show entry
$ gopass show entry key
$ gopass show entry --qr
$ gopass show entry --password
```

## Modes of operation

* Show the whole entry: `gopass show entry`
* Show a specific key of the given entry: `gopass show entry key` (only works for key-value or YAML secrets)

## Flags

Flag | Aliases | Description
---- | ------- | -----------
`--clip` | `-c` | Copy the password value into the clipboard and don't show the content.
`--qr` | | Encode the password field as a QR code and print it. Note: When combining with `-c` the unencoded password is copied, not the QR code.
`--password` | `-o` | Display only the password. For use in scripts. Takes precedence over other flags.
`--revision` | `-r` | Display a specific revision of the entry. Use an exact version identifier from `gopass history` or the special `-<N>` syntax. Does not work with native (e.g. git) refs.
`--noparsing` | `-n` | Do not parse the content, disable YAML and Key-Value functions.
`--chars` | | Display selected characters from the password.

## Details

This section describes the expected behaviour of the `show` command with respect to different combinations of flags and
config options.

Note: This section describes the expected behaviour, not necessarily the observed behaviour.
If you notice any discrepancies please file a bug and we will try to fix it.

TODO: We need to specify the expectations around new lines.

* When no flag is set the `show` command will display the full content of the secret and will parse it to support key-value lookup and YAML entries.
* The `--noparsing` flag will disable all parsing of the output, this can help debugging YAML secrets for example, where `key: 0123` actually parses into octal for 83. 
* The `--clip` flag will copy the value of the `Password` field to the clipboard and doesn't display any part of the secret.
* The `--qr` flags operates complementary to other flags. It will *additionally* format the value of the `Password` entry as a QR code and display it. Other than that it will honor the other options, e.g. `gopass show --qr` will display the QR code *and* the whole secret content below. One special case is the `-o` flag, this flag doesn't make a lot of sense in combination, so if both `--qr` and `-o` are given only the QR code will be displayed.
* Arbitrary git refs are not supported as arguments to the `--revision` flag. Using those might work, but this is explicitly not supported and bug reports will be closed as `wont-fix`. The main issue with using arbitrary git refs is that git versions a whole repository, not single files. So the revision `HEAD^` might not have any changes for a given entry. Thus we only support specifc revisions obtained from `gopass history` or our custom syntax `-N` where N is an integer identifying a specific commit before `HEAD` (cf. `HEAD~N`).

## Parsing and secrets

Secrets are stored on disk as provided, but are parsed upon display to provide extra features such as the ability 
to show the value of a key using:  `gopass show entry key`.

The secrets are split into 3 categories:
 - the plain type, which is just a plain secret without key-value capabilities 
    ```
    this is a plain secret
    using multiple lines
    
    and that's it
    ```
    gets parsed to the same value


 - the key-value type, which allows to query the value of a specific key. This does not preserve ordering.
    ```
    this is a KV secret
    where: the first line is the password
    and: the keys are separated from their value by :
    
    and maybe we have a body text
    below it
    ```
    will be parsed into:
   ```
    this is a KV secret
    and: the keys are separated from their value by :
    where: the first line is the password
    
    
    and maybe we have a body text
    below it
    ```


 - the YAML type which implements YAML support, which means that secrets are parsed as per YAML standard.
    ```
    s3cret
    ---
    invoice: 0123
    date   : 2001-01-23
    bill-to: &id001
        given  : Bob
        family : Doe
    ship-to: *id001
    ```
   will be parsed into:
    ```
    s3cret
    ---
    bill-to:
        family: Doe
        given: Bob
    date: 2001-01-23T00:00:00Z
    invoice: 83
    ship-to:
        family: Doe
        given: Bob
    ```
   Note how the `0123` is interpreted as octal for 83. If you want to store a string made of digits such as a numerical
   username, it should be enclosed in string delimiters: `username: "0123"` will always be parsed as the string `0123`
   and not as octal.
