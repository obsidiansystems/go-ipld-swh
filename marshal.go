package ipldswh

import (
	"fmt"
	"io"
)

// Encode serializes a git node to a raw binary form.
func Encode(n Snapshot, w io.Writer) error {
	if _, err := fmt.Fprintf(w, "snapshot "); err != nil {
		return err
	}
	return EncodeSnapshot(n, w)
}
