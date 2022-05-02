package ipldswh

//go:generate go run ./gen .
//go:generate go fmt ./

import (
	"io"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	mc "github.com/ipld/go-ipld-prime/multicodec"
)

var (
	d ipld.Decoder = func(na ipld.NodeAssembler, r io.Reader) error {
		x, err := Decode(r)
		if err != nil {
			return err
		}
		na.AssignNode(x)
		return nil
	}
	e ipld.Encoder = func(n ipld.Node, w io.Writer) error {
		builder := Type.Snapshot.NewBuilder()
		if err := builder.AssignNode(n); err != nil {
			return err
		}
		return Encode(builder.Build().(Snapshot), w)
	}
)

func init() {
	mc.RegisterEncoder(cid.GitRaw, e)
	mc.RegisterDecoder(cid.GitRaw, d)
}
