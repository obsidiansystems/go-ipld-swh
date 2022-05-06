package ipldswh

//go:generate go run ./gen .
//go:generate go fmt ./

import (
	"github.com/ipld/go-ipld-prime"
	mc "github.com/ipld/go-ipld-prime/multicodec"
)

var (
	_ ipld.Encoder = EncodeGeneric
	_ ipld.Decoder = DecodeGeneric
)

func init() {
	mc.RegisterEncoder(Swh1Snp, EncodeGeneric)
	mc.RegisterDecoder(Swh1Snp, DecodeGeneric)
}
