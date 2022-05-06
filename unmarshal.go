package ipldswh

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ipld/go-ipld-prime"
)

// Decode produces an Snapshot from a stream / binary represnetation.
// Decode reads from a reader to fill a NodeAssembler
func Decode(r io.Reader) (Snapshot, error) {
	rd := bufio.NewReader(r)

	typ, err := rd.ReadString(' ')
	if err == io.EOF {
		return nil, io.ErrUnexpectedEOF
	}
	if err != nil {
		return nil, err
	}
	typ = typ[:len(typ)-1]

	switch typ {
	case "snapshot":
		return DecodeSnapshot(rd)
	default:
		return nil, fmt.Errorf("unrecognized object type: %q", typ)
	}
}

func DecodeGeneric(na ipld.NodeAssembler, r io.Reader) error {
	x, err := Decode(r)
	if err != nil {
		return err
	}
	na.AssignNode(x.Representation())
	return nil
}
