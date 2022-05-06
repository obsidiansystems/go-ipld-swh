package ipldswh

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/schema"
)

// DecodeSnapshot fills a NodeAssembler (from `Type.Snapshot__Repr.NewBuilder()`) from a stream of bytes
func DecodeSnapshot(rd *bufio.Reader) (Snapshot, error) {
	if _, err := readTerminatedNumber(rd, 0); err != nil {
		return nil, err
	}

	var r _Snapshot
	r.m = make(map[_BranchName]MaybeSnapshotBranch)
	r.t = make([]_Snapshot__entry, 0)

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
	value []byte
}
type namedBranchRaw struct {
	key string
	b   branchRaw
}

func decodeSnapshotBranchRaw(rd *bufio.Reader) (*namedBranchRaw, error) {
	ty, err := rd.ReadString(' ')
	if err != nil {
		return nil, err
	}
	ty = ty[:len(ty)-1]

	key, err := rd.ReadString(0)
	if err != nil {
		return nil, err
	}
	key = key[:len(key)-1]

	valueLen, err := readTerminatedNumber(rd, ':')
	if err != nil {
		return nil, err
	}
	value := make([]byte, valueLen)
	_, err = io.ReadFull(rd, value)
	if err != nil {
		return nil, err
	}

	return &namedBranchRaw{key: key, b: branchRaw{ty, value}}, nil
}

// DecodeSnapshotBranch fills a NodeAssembler (from `Type.SnapshotBranch__Repr.NewBuilder()`) from a stream of bytes
func DecodeSnapshotBranch(rd *bufio.Reader) (string, MaybeSnapshotBranch, error) {
	raw, err := decodeSnapshotBranchRaw(rd)
	if err != nil {
		return "", nil, err
	}

	gitLink := func(value []byte) _Link {
		return _Link{x: cidlink.Link{Cid: gitShaToCid(value)}}
	}

	swh1SnpLink := func(value []byte) _Snapshot_Link {
		return _Snapshot_Link{x: cidlink.Link{Cid: swh1SnpShaToCid(value)}}
	}

	var te MaybeSnapshotBranch = nil
	switch raw.b.ty {
	case "dangling":
		if len(raw.b.value) != 0 {
			return "",
				nil,
				fmt.Errorf("Dangling branch should not have non-empty value %q", raw.b.value)
		}
		te = &_Snapshot__valueAbsent
	case "content":
		te = &_SnapshotBranch__Maybe{
			m: schema.Maybe_Value,
			v: &_SnapshotBranch{
				tag: 1,
				x1: _SnapshotBranch_Content{
					target: gitLink(raw.b.value),
				},
			},
		}
	case "directory":
		te = &_SnapshotBranch__Maybe{
			m: schema.Maybe_Value,
			v: &_SnapshotBranch{
				tag: 2,
				x2: _SnapshotBranch_Directory{
					target: gitLink(raw.b.value),
				},
			},
		}
	case "revision":
		te = &_SnapshotBranch__Maybe{
			m: schema.Maybe_Value,
			v: &_SnapshotBranch{
				tag: 3,
				x3: _SnapshotBranch_Revision{
					target: gitLink(raw.b.value),
				},
			},
		}
	case "release":
		te = &_SnapshotBranch__Maybe{
			m: schema.Maybe_Value,
			v: &_SnapshotBranch{
				tag: 4,
				x4: _SnapshotBranch_Release{
					target: gitLink(raw.b.value),
				},
			},
		}
	case "snapshot":
		te = &_SnapshotBranch__Maybe{
			m: schema.Maybe_Value,
			v: &_SnapshotBranch{
				tag: 5,
				x5: _SnapshotBranch_Snapshot{
					target: swh1SnpLink(raw.b.value),
				},
			},
		}
	case "alias":
		te = &_SnapshotBranch__Maybe{
			m: schema.Maybe_Value,
			v: &_SnapshotBranch{
				tag: 6,
				x6: _SnapshotBranch_Alias{
					target: _BranchName{x: string(raw.b.value)},
				},
			},
		}
	default:
		return "",
			nil,
			fmt.Errorf("unrecognized snapshot branch type: %q", raw.b.ty)
	}
	return raw.key, te, nil
}

func EncodeSnapshot(s Snapshot, w io.Writer) error {
	buf := new(bytes.Buffer)

	for key, te := range s.m {
		raw0, err := encodeSnapshotBranchRaw(te)
		if err != nil {
			return err
		}
		raw := namedBranchRaw{key: key.x, b: *raw0}
		if err := EncodeSnapshotBranch(raw, buf); err != nil {
			return err
		}
	}
	cnt := buf.Len()
	if _, err := fmt.Fprintf(w, "%d\x00", cnt); err != nil {
		return err
	}

	_, err := buf.WriteTo(w)
	return err
}

func encodeSnapshotBranchRaw(n MaybeSnapshotBranch) (*branchRaw, error) {

	castLink := func(l ipld.Link) (*cidlink.Link, error) {
		cl, ok := l.(cidlink.Link)
		if !ok {
			// this _should_ be taken care of by the Typed conversion above with
			// "missing required fields: Hash"
			return nil, fmt.Errorf("invalid DAG-PB form (link must have a Hash)")
		}
		return &cl, nil
	}
	switch n.m {
	case schema.Maybe_Null:
		return &branchRaw{
			ty:    "absent",
			value: []byte{},
		}, nil
	case schema.Maybe_Value:
		n1 := n.v
		switch n1.tag {
		case 1:
			n2 := n1.x1
			cl, err := castLink(n2.target.x)
			if err != nil {
				return nil, err
			}
			return &branchRaw{
				ty:    "content",
				value: cidToSha(cl.Cid),
			}, nil
		case 2:
			n2 := n1.x2
			cl, err := castLink(n2.target.x)
			if err != nil {
				return nil, err
			}
			return &branchRaw{
				ty:    "directory",
				value: cidToSha(cl.Cid),
			}, nil
		case 3:
			n2 := n1.x3
			cl, err := castLink(n2.target.x)
			if err != nil {
				return nil, err
			}
			return &branchRaw{
				ty:    "revision",
				value: cidToSha(cl.Cid),
			}, nil
		case 4:
			n2 := n1.x4
			cl, err := castLink(n2.target.x)
			if err != nil {
				return nil, err
			}
			return &branchRaw{
				ty:    "release",
				value: cidToSha(cl.Cid),
			}, nil
		case 5:
			n2 := n1.x5
			cl, err := castLink(n2.target.x)
			if err != nil {
				return nil, err
			}
			return &branchRaw{
				ty:    "snapshot",
				value: cidToSha(cl.Cid),
			}, nil
		default:
			panic("invalid union state; how did you create this object?")
		}
	default:
		panic("invalid maybe state; how did you create this object?")
	}
}

func EncodeSnapshotBranch(n namedBranchRaw, w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s %s\x00%d:", n.b.ty, n.key, len(n.b.value))
	if err != nil {
		return err
	}

	_, err = w.Write(n.b.value)
	if err != nil {
		return err
	}

	return nil
}
