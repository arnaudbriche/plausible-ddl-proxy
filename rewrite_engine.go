package main

import (
	"fmt"
	"regexp"
	"strings"
)

var bypassEngineRe = regexp.MustCompile(`ReplicatedVersionedCollapsingMergeTree|ReplicatedMergeTree|ReplicatedCollapsingMergeTree|ReplicatedSummingMergeTree`)

func bypassEngine(rctx rewriteContext) func(string) string {
	return func(s string) string { return s }
}

var tinyLogRe = regexp.MustCompile(`TinyLog`)

func rewriteTinyLog(rctx rewriteContext) func(string) string {
	return func(s string) string {
		return fmt.Sprintf(
			"ReplicatedMergeTree('%s', '%s') ORDER BY tuple()",
			rctx.ZkPathPrefix+rctx.TableName,
			rctx.Replica,
		)
	}
}

var mergeTreeRe = regexp.MustCompile(`MergeTree[\(\)]*`)

func rewriteMergeTree(rctx rewriteContext) func(string) string {
	return func(s string) string {
		return fmt.Sprintf(
			"ReplicatedMergeTree('%s', '%s')",
			rctx.ZkPathPrefix+rctx.TableName,
			rctx.Replica,
		)
	}
}

var collapsingMergeTreeRe = regexp.MustCompile(`CollapsingMergeTree\(`)

func rewriteCollapsingMergeTree(rctx rewriteContext) func(string) string {
	return func(s string) string {
		return fmt.Sprintf(
			"ReplicatedCollapsingMergeTree('%s', '%s', %s",
			rctx.ZkPathPrefix+rctx.TableName,
			rctx.Replica,
			strings.ReplaceAll(s, "CollapsingMergeTree(", ""),
		)
	}
}

var summingMergeTreeRe = regexp.MustCompile(`SummingMergeTree\(`)

func rewriteSummingMergeTreeRe(rctx rewriteContext) func(string) string {
	return func(s string) string {
		return fmt.Sprintf(
			"ReplicatedSummingMergeTree('%s', '%s', %s",
			rctx.ZkPathPrefix+rctx.TableName,
			rctx.Replica,
			strings.ReplaceAll(s, "SummingMergeTree(", ""),
		)
	}
}

var versionedCollapsingMergeTreeRe = regexp.MustCompile(`VersionedCollapsingMergeTree\(`)

func rewriteVersionedCollapsingMergeTree(rctx rewriteContext) func(string) string {
	return func(s string) string {
		return fmt.Sprintf(
			"ReplicatedVersionedCollapsingMergeTree('%s', '%s', %s",
			rctx.ZkPathPrefix+rctx.TableName,
			rctx.Replica,
			strings.ReplaceAll(s, "VersionedCollapsingMergeTree(", ""),
		)
	}
}

var engineRewrites = rewrites{
	{Regex: bypassEngineRe, Callback: bypassEngine},
	{Regex: versionedCollapsingMergeTreeRe, Callback: rewriteVersionedCollapsingMergeTree},
	{Regex: summingMergeTreeRe, Callback: rewriteSummingMergeTreeRe},
	{Regex: collapsingMergeTreeRe, Callback: rewriteCollapsingMergeTree},
	{Regex: mergeTreeRe, Callback: rewriteMergeTree},
	{Regex: tinyLogRe, Callback: rewriteTinyLog},
}
