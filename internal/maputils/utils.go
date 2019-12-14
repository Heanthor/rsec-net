package maputils

import (
	"bytes"
	"crypto/md5"
	"fmt"

	cmap "github.com/orcaman/concurrent-map"
)

// ComputeHash computes the md5 hash of the given cmap.
// This is will lock the entire map, but shouldn't be too slow
func ComputeHash(cmap cmap.ConcurrentMap) (md5Hash [16]byte, items map[string]interface{}) {
	var b bytes.Buffer
	items = cmap.Items()

	for k, v := range items {
		fmt.Fprintf(&b, "%s=%+v", k, v)
	}
	md5Hash = md5.Sum(b.Bytes())

	return
}
