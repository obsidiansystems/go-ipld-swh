package ipldswh

//go:generate go run ./gen .
//go:generate go fmt ./

import (
	"github.com/ipld/go-ipld-prime"
	mc2 "github.com/ipld/go-ipld-prime/multicodec"
	mc "github.com/multiformats/go-multicodec"
)

var (
	_ ipld.Encoder = EncodeGeneric
	_ ipld.Decoder = DecodeGeneric
)

func init() {
	mc2.RegisterEncoder((uint64)(mc.Swhid1Snp), EncodeGeneric)
	mc2.RegisterDecoder((uint64)(mc.Swhid1Snp), DecodeGeneric)
}
