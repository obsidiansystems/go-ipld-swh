package ipldswh

import (
	"fmt"
	"io"

	"github.com/ipld/go-ipld-prime"
)

// Encode serializes a git node to a raw binary form.
func Encode(n ipld.Node, w io.Writer) error {
	switch n.Prototype() {
	case Type.Snapshot, Type.Snapshot__Repr:
		return encodeSnapshot(n, w)
	default:
	}
	switch n.Kind() {
	case ipld.Kind_List:
		return encodeSnapshot(n, w)
	default:
		return fmt.Errorf("unrecognized object type: %T", n.Prototype())
	}
}
