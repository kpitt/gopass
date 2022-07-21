package store

import "fmt"

var (
	// ErrExistsFailed is returend if we can't check for existence.
	ErrExistsFailed = fmt.Errorf("Failed to check for existence")
	// ErrNotFound is returned if an entry was not found.
	ErrNotFound = fmt.Errorf("Entry is not in the password store")
	// ErrEncrypt is returned if we failed to encrypt an entry.
	ErrEncrypt = fmt.Errorf("Failed to encrypt")
	// ErrDecrypt is returned if we failed to decrypt and entry.
	ErrDecrypt = fmt.Errorf("Failed to decrypt")
	// ErrIO is any kind of I/O error.
	ErrIO = fmt.Errorf("I/O error")
	// ErrGitInit is returned if git is already initialized.
	ErrGitInit = fmt.Errorf("git is already initialized")
	// ErrGitNotInit is returned if git is not initialized.
	ErrGitNotInit = fmt.Errorf("git is not initialized")
	// ErrGitNoRemote is returned if git has no origin remote.
	ErrGitNoRemote = fmt.Errorf("git has no remote origin")
	// ErrGitNothingToCommit is returned if there are no staged changes.
	ErrGitNothingToCommit = fmt.Errorf("git has nothing to commit")
	// ErrEmptySecret is returned if a secret exists but has no content.
	ErrEmptySecret = fmt.Errorf("Empty secret")
	// ErrNoBody is returned if a secret exists but has no content beyond a password.
	ErrNoBody = fmt.Errorf("No safe content to display. Use -f to force display.")
	// ErrNoPassword is returned is a secret exists but has no password, only a body.
	ErrNoPassword = fmt.Errorf("No password to display. Check the body of the entry instead.")
	// ErrYAMLNoMark is returned if a secret contains no valid YAML document marker.
	ErrYAMLNoMark = fmt.Errorf("No YAML document marker found")
	// ErrNoKey is returned if a KV or YAML entry doesn't contain a key.
	ErrNoKey = fmt.Errorf("Key not found in entry")
	// ErrYAMLValueUnsupported is returned if the user tries to unmarshal an nested struct.
	ErrYAMLValueUnsupported = fmt.Errorf("Cannot unmarshal nested YAML value")
)
