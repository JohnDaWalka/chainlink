package crypto

import "strings"

type CSAKey struct {
	Key string
}

func (c *CSAKey) CleansedKey() string {
	return strings.TrimPrefix(c.Key, "csa_")
}
