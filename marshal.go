package ipldswh

import (
	"fmt"
	"io"

	"github.com/ipld/go-ipld-prime"
)

// Encode serializes a git node to a raw binary form.
func Encode(n Snapshot, w io.Writer) error {
	if _, err := fmt.Fprintf(w, "snapshot "); err != nil {
		return err
	}
	return EncodeSnapshot(n, w)
}

func EncodeGeneric(n ipld.Node, w io.Writer) error {
	var assembler _Snapshot__ReprAssembler
	if err := assembler.AssignNode(n); err != nil {
		return err
	}
	return Encode(assembler.w, w)
}
