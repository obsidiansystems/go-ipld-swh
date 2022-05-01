package ipldswh

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/schema"
)

// DecodeSnapshot fills a NodeAssembler (from `Type.Snapshot__Repr.NewBuilder()`) from a stream of bytes
func DecodeSnapshot(rd *bufio.Reader) (Snapshot, error) {
	if _, err := readTerminatedNumber(rd, 0); err != nil {
		return nil, err
	}

	var r _Snapshot
	for {
		if _, err := rd.Peek(1); err == io.EOF {
			break
		}
		name, mbranch, err := DecodeSnapshotBranch(rd)
		if err != nil {
			return nil, err
		}
		k := _BranchName{x: name}
		r.m[k] = mbranch
		r.t = append(r.t, _Snapshot__entry{k: k, v: *mbranch})
	}
	return &r, nil
}

type branchRaw struct {
	ty    string
	key   string
	value []byte
}

func decodeSnapshotBranchRaw(rd *bufio.Reader) (branchRaw, error) {
	ty, err := rd.ReadString(' ')
	if err != nil {
		return branchRaw{}, err
	}
	ty = ty[:len(ty)-1]

	key, err := rd.ReadString(0)
	if err != nil {
		return branchRaw{}, err
	}
	key = key[:len(key)-1]

	valueLen, err := readTerminatedNumber(rd, ':')
	if err != nil {
		return branchRaw{}, err
	}
	value := make([]byte, valueLen)
	_, err = io.ReadFull(rd, value)
	if err != nil {
		return branchRaw{}, err
	}

	return branchRaw{ty, key, value}, nil
}

// DecodeSnapshotBranch fills a NodeAssembler (from `Type.SnapshotBranch__Repr.NewBuilder()`) from a stream of bytes
func DecodeSnapshotBranch(rd *bufio.Reader) (string, MaybeSnapshotBranch, error) {
	raw, err := decodeSnapshotBranchRaw(rd)
	if err != nil {
		return "", nil, err
	}

	gitLink := func(value []byte) _Link {
		return _Link{cidlink.Link{Cid: gitShaToCid(value)}}
	}

	swh1SnpLink := func(value []byte) _Snapshot_Link {
		return _Snapshot_Link{cidlink.Link{Cid: swh1SnpShaToCid(value)}}
	}

	var te MaybeSnapshotBranch = nil
	switch raw.ty {
	case "dangling":
		if len(raw.value) != 0 {
			return "",
				nil,
				fmt.Errorf("Dangling branch should not have non-empty value %q", raw.value)
		}
		te = &_Snapshot__valueAbsent
	case "content":
		te = &_SnapshotBranch__Maybe{
			m: schema.Maybe_Value,
			v: &_SnapshotBranch{
				tag: 1,
				x1: _SnapshotBranch_Content{
					target: gitLink(raw.value),
				},
			},
		}
	case "directory":
		te = &_SnapshotBranch__Maybe{
			m: schema.Maybe_Value,
			v: &_SnapshotBranch{
				tag: 2,
				x2: _SnapshotBranch_Directory{
					target: gitLink(raw.value),
				},
			},
		}
	case "revision":
		te = &_SnapshotBranch__Maybe{
			m: schema.Maybe_Value,
			v: &_SnapshotBranch{
				tag: 3,
				x3: _SnapshotBranch_Revision{
					target: gitLink(raw.value),
				},
			},
		}
	case "release":
		te = &_SnapshotBranch__Maybe{
			m: schema.Maybe_Value,
			v: &_SnapshotBranch{
				tag: 4,
				x4: _SnapshotBranch_Release{
					target: gitLink(raw.value),
				},
			},
		}
	case "snapshot":
		te = &_SnapshotBranch__Maybe{
			m: schema.Maybe_Value,
			v: &_SnapshotBranch{
				tag: 5,
				x5: _SnapshotBranch_Snapshot{
					target: swh1SnpLink(raw.value),
				},
			},
		}
	case "alias":
		te = &_SnapshotBranch__Maybe{
			m: schema.Maybe_Value,
			v: &_SnapshotBranch{
				tag: 6,
				x6: _SnapshotBranch_Alias{
					target: _BranchName{x: string(raw.value)},
				},
			},
		}
	default:
		return "",
			nil,
			fmt.Errorf("unrecognized snapshot branch type: %q", raw.ty)
	}
	return raw.key, te, nil
}

func encodeSnapshot(n ipld.Node, w io.Writer) error {
	buf := new(bytes.Buffer)

	mi := n.MapIterator()
	for !mi.Done() {
		key, te, err := mi.Next()
		if err != nil {
			return err
		}
		name, err := key.AsString()
		if err != nil {
			return err
		}
		if err := encodeSnapshotBranch(name, te, buf); err != nil {
			return err
		}
	}
	cnt := buf.Len()
	if _, err := fmt.Fprintf(w, "snapshot %d\x00", cnt); err != nil {
		return err
	}

	_, err := buf.WriteTo(w)
	return err
}

func encodeSnapshotBranch(name string, n ipld.Node, w io.Writer) error {
	i := n.MapIterator()
	ty, n1, err := i.Next()
	i.Done()

	var v []byte

	tyS, err := ty.AsString()
	if err != nil {
		return err
	}

	switch tyS {
	case "content":
		n3, err := n1.LookupByString("target")
		if err != nil {
			return err
		}
		l, err := n3.AsLink()
		if err != nil {
			return err
		}
		v = cidToSha(l.(cidlink.Link).Cid)
	case "directory":
		n3, err := n1.LookupByString("target")
		if err != nil {
			return err
		}
		l, err := n3.AsLink()
		if err != nil {
			return err
		}
		v = cidToSha(l.(cidlink.Link).Cid)
	case "revision":
		n3, err := n1.LookupByString("target")
		if err != nil {
			return err
		}
		l, err := n3.AsLink()
		if err != nil {
			return err
		}
		v = cidToSha(l.(cidlink.Link).Cid)
	case "release":
		n3, err := n1.LookupByString("target")
		if err != nil {
			return err
		}
		l, err := n3.AsLink()
		if err != nil {
			return err
		}
		v = cidToSha(l.(cidlink.Link).Cid)
	case "snapshot":
		n3, err := n1.LookupByString("target")
		if err != nil {
			return err
		}
		l, err := n3.AsLink()
		if err != nil {
			return err
		}
		v = cidToSha(l.(cidlink.Link).Cid)
	case "alias":
		n3, err := n1.LookupByString("target")
		if err != nil {
			return err
		}
		s, err := n3.AsString()
		if err != nil {
			return err
		}
		v = []byte(s)
	case "dangling":
		n3, err := n1.LookupByString("target")
		if err != nil {
			return err
		}
		b, err := n3.AsBytes()
		if err != nil {
			return err
		}
		if len(b) != 0 {
			return fmt.Errorf("Dangling branch should not have non-empty value %q", b)
		}
	default:
		return err
	}

	_, err = fmt.Fprintf(w, "%s %s\x00%d:", ty, name, len(v))
	if err != nil {
		return err
	}

	_, err = w.Write(v)
	if err != nil {
		return err
	}

	return nil
}
