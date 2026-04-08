package experiment

import (
	"crypto/sha256"
	"encoding/binary"
	"sort"
)

// VariantForUser returns a sticky variant index for (userKey, experimentKey).
func VariantForUser(userID string, experimentKey string, variants []string) string {
	if len(variants) == 0 {
		return ""
	}
	vs := append([]string(nil), variants...)
	sort.Strings(vs)
	h := sha256.Sum256([]byte(userID + "\x00" + experimentKey))
	idx := binary.BigEndian.Uint64(h[:8]) % uint64(len(vs))
	return vs[idx]
}
