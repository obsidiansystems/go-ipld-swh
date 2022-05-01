package ipldswh

import (
	"bufio"
	"fmt"
	"io"
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
