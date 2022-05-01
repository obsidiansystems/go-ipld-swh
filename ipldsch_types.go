package ipldswh

// Code generated by go-ipld-prime gengo.  DO NOT EDIT.

import (
	ipld "github.com/ipld/go-ipld-prime"
)

var _ ipld.Node = nil // suppress errors when this dependency is not referenced
// Type is a struct embeding a NodePrototype/Type for every Node implementation in this package.
// One of its major uses is to start the construction of a value.
// You can use it like this:
//
// 		ipldswh.Type.YourTypeName.NewBuilder().BeginMap() //...
//
// and:
//
// 		ipldswh.Type.OtherTypeName.NewBuilder().AssignString("x") // ...
//
var Type typeSlab

type typeSlab struct {
	BranchName                     _BranchName__Prototype
	BranchName__Repr               _BranchName__ReprPrototype
	Link                           _Link__Prototype
	Link__Repr                     _Link__ReprPrototype
	Snapshot                       _Snapshot__Prototype
	Snapshot__Repr                 _Snapshot__ReprPrototype
	SnapshotBranch                 _SnapshotBranch__Prototype
	SnapshotBranch__Repr           _SnapshotBranch__ReprPrototype
	SnapshotBranch_Alias           _SnapshotBranch_Alias__Prototype
	SnapshotBranch_Alias__Repr     _SnapshotBranch_Alias__ReprPrototype
	SnapshotBranch_Content         _SnapshotBranch_Content__Prototype
	SnapshotBranch_Content__Repr   _SnapshotBranch_Content__ReprPrototype
	SnapshotBranch_Directory       _SnapshotBranch_Directory__Prototype
	SnapshotBranch_Directory__Repr _SnapshotBranch_Directory__ReprPrototype
	SnapshotBranch_Release         _SnapshotBranch_Release__Prototype
	SnapshotBranch_Release__Repr   _SnapshotBranch_Release__ReprPrototype
	SnapshotBranch_Revision        _SnapshotBranch_Revision__Prototype
	SnapshotBranch_Revision__Repr  _SnapshotBranch_Revision__ReprPrototype
	SnapshotBranch_Snapshot        _SnapshotBranch_Snapshot__Prototype
	SnapshotBranch_Snapshot__Repr  _SnapshotBranch_Snapshot__ReprPrototype
	Snapshot_Link                  _Snapshot_Link__Prototype
	Snapshot_Link__Repr            _Snapshot_Link__ReprPrototype
	String                         _String__Prototype
	String__Repr                   _String__ReprPrototype
}

// --- type definitions follow ---

// BranchName matches the IPLD Schema type "BranchName".  It has string kind.
type BranchName = *_BranchName
type _BranchName struct{ x string }

// Link matches the IPLD Schema type "Link".  It has link kind.
type Link = *_Link
type _Link struct{ x ipld.Link }

// Snapshot matches the IPLD Schema type "Snapshot".  It has map kind.
type Snapshot = *_Snapshot
type _Snapshot struct {
	m map[_BranchName]MaybeSnapshotBranch
	t []_Snapshot__entry
}
type _Snapshot__entry struct {
	k _BranchName
	v _SnapshotBranch__Maybe
}

// SnapshotBranch matches the IPLD Schema type "SnapshotBranch".
// SnapshotBranch has Union typekind, which means its data model behaviors are that of a map kind.
type SnapshotBranch = *_SnapshotBranch
type _SnapshotBranch struct {
	tag uint
	x1  _SnapshotBranch_Content
	x2  _SnapshotBranch_Directory
	x3  _SnapshotBranch_Revision
	x4  _SnapshotBranch_Release
	x5  _SnapshotBranch_Snapshot
	x6  _SnapshotBranch_Alias
}
type _SnapshotBranch__iface interface {
	_SnapshotBranch__member()
}

func (_SnapshotBranch_Content) _SnapshotBranch__member()   {}
func (_SnapshotBranch_Directory) _SnapshotBranch__member() {}
func (_SnapshotBranch_Revision) _SnapshotBranch__member()  {}
func (_SnapshotBranch_Release) _SnapshotBranch__member()   {}
func (_SnapshotBranch_Snapshot) _SnapshotBranch__member()  {}
func (_SnapshotBranch_Alias) _SnapshotBranch__member()     {}

// SnapshotBranch_Alias matches the IPLD Schema type "SnapshotBranch_Alias".  It has Struct type-kind, and may be interrogated like map kind.
type SnapshotBranch_Alias = *_SnapshotBranch_Alias
type _SnapshotBranch_Alias struct {
	target _BranchName
}

// SnapshotBranch_Content matches the IPLD Schema type "SnapshotBranch_Content".  It has Struct type-kind, and may be interrogated like map kind.
type SnapshotBranch_Content = *_SnapshotBranch_Content
type _SnapshotBranch_Content struct {
	target _Link
}

// SnapshotBranch_Directory matches the IPLD Schema type "SnapshotBranch_Directory".  It has Struct type-kind, and may be interrogated like map kind.
type SnapshotBranch_Directory = *_SnapshotBranch_Directory
type _SnapshotBranch_Directory struct {
	target _Link
}

// SnapshotBranch_Release matches the IPLD Schema type "SnapshotBranch_Release".  It has Struct type-kind, and may be interrogated like map kind.
type SnapshotBranch_Release = *_SnapshotBranch_Release
type _SnapshotBranch_Release struct {
	target _Link
}

// SnapshotBranch_Revision matches the IPLD Schema type "SnapshotBranch_Revision".  It has Struct type-kind, and may be interrogated like map kind.
type SnapshotBranch_Revision = *_SnapshotBranch_Revision
type _SnapshotBranch_Revision struct {
	target _Link
}

// SnapshotBranch_Snapshot matches the IPLD Schema type "SnapshotBranch_Snapshot".  It has Struct type-kind, and may be interrogated like map kind.
type SnapshotBranch_Snapshot = *_SnapshotBranch_Snapshot
type _SnapshotBranch_Snapshot struct {
	target _Snapshot_Link
}

// Snapshot_Link matches the IPLD Schema type "Snapshot_Link".  It has link kind.
type Snapshot_Link = *_Snapshot_Link
type _Snapshot_Link struct{ x ipld.Link }

// String matches the IPLD Schema type "String".  It has string kind.
type String = *_String
type _String struct{ x string }
