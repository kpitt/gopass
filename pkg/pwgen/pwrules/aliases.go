package pwrules

import (
	"sort"
)

// LookupAliases looks up known aliases for the given domain.
func LookupAliases(domain string) []string {
	// TODO: Is there a more efficient way to copy a single array?
	aliases := make([]string, 0, len(genAliases[domain]))
	aliases = append(aliases, genAliases[domain]...)
	sort.Strings(aliases)

	return aliases
}

// AllAliases returns all aliases.
func AllAliases() map[string][]string {
	all := make(map[string][]string, len(genAliases))
	for k, v := range genAliases {
		all[k] = append(all[k], v...)
	}

	return all
}
