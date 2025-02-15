//go:build openbsd
// +build openbsd

package protect

import "golang.org/x/sys/unix"

// Pledge on OpenBSD lets us "promise" to only run a subset of
// system calls: http://man.openbsd.org/pledge
func Pledge(s string) error {
	return unix.PledgePromises(s)
}
