package ipldswh

import (
	"io"
)

// Encode serializes a git node to a raw binary form.
func Encode(n Snapshot, w io.Writer) error {
	return encodeSnapshot(n, w)
}
