package main

import (
	"fmt"
	"os"

	"github.com/ipld/go-ipld-prime/schema"
	gengo "github.com/ipld/go-ipld-prime/schema/gen/go"
)

func main() {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &gengo.AdjunctCfg{
		CfgUnionMemlayout: map[schema.TypeName]string{},
		FieldSymbolLowerOverrides: map[gengo.FieldTuple]string{
			{TypeName: "Tag", FieldName: "type"}: "typ",
		},
	}

	ts.Accumulate(schema.SpawnString("String")) // Needed for Union Tag
	ts.Accumulate(schema.SpawnString("BranchName"))
	ts.Accumulate(schema.SpawnLink("Link"))

	ts.Accumulate(schema.SpawnStruct("SnapshotBranch_Content", []schema.StructField{
		schema.SpawnStructField("target", "Link", false, false),
	}, schema.SpawnStructRepresentationMap(map[string]string{})))

	ts.Accumulate(schema.SpawnStruct("SnapshotBranch_Directory", []schema.StructField{
		schema.SpawnStructField("target", "Link", false, false),
	}, schema.SpawnStructRepresentationMap(map[string]string{})))

	ts.Accumulate(schema.SpawnStruct("SnapshotBranch_Revision", []schema.StructField{
		schema.SpawnStructField("target", "Link", false, false),
	}, schema.SpawnStructRepresentationMap(map[string]string{})))

	ts.Accumulate(schema.SpawnStruct("SnapshotBranch_Release", []schema.StructField{
		schema.SpawnStructField("target", "Link", false, false),
	}, schema.SpawnStructRepresentationMap(map[string]string{})))

	ts.Accumulate(schema.SpawnStruct("SnapshotBranch_Snapshot", []schema.StructField{
		schema.SpawnStructField("target", "Snapshot_Link", false, false),
	}, schema.SpawnStructRepresentationMap(map[string]string{})))

	ts.Accumulate(schema.SpawnStruct("SnapshotBranch_Alias", []schema.StructField{
		schema.SpawnStructField("target", "BranchName", false, false),
	}, schema.SpawnStructRepresentationMap(map[string]string{})))

	ts.Accumulate(schema.SpawnUnion("SnapshotBranch", []schema.TypeName{
		"SnapshotBranch_Content",
		"SnapshotBranch_Directory",
		"SnapshotBranch_Revision",
		"SnapshotBranch_Release",
		"SnapshotBranch_Snapshot",
		"SnapshotBranch_Alias",
	}, schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
		"content": "SnapshotBranch_Content",
		"directory": "SnapshotBranch_Directory",
		"revision": "SnapshotBranch_Revision",
		"release": "SnapshotBranch_Release",
		"snapshot": "SnapshotBranch_Snapshot",
		"alias": "SnapshotBranch_Alias",
	})))

	ts.Accumulate(schema.SpawnMap("Snapshot", "BranchName", "SnapshotBranch", true))

	ts.Accumulate(schema.SpawnLinkReference("Snapshot_Link", "Snapshot"))

	if errs := ts.ValidateGraph(); errs != nil {
		for _, err := range errs {
			fmt.Printf("- %s\n", err)
		}
		panic("not happening")
	}

	gengo.Generate(os.Args[1], "ipldswh", ts, adjCfg)
}
