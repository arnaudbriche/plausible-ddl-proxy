package main

import (
	"fmt"
	"regexp"
)

var bypassOnClusterRe = regexp.MustCompile(`ON CLUSTER`)

func bypassOnCluster(rctx rewriteContext) func(string) string {
	return func(s string) string { return s }
}

var createTableIfNotExistsRe = regexp.MustCompile(`CREATE TABLE IF NOT EXISTS [\"\w]+`)

func rewriteCreateTableIfNotExists(rctx rewriteContext) func(string) string {
	return func(s string) string {
		return fmt.Sprintf(
			"%s ON CLUSTER '%s'",
			s,
			rctx.Cluster,
		)
	}
}

var onClusterRewrites = rewrites{
	{Regex: bypassOnClusterRe, Callback: bypassOnCluster},
	{Regex: createTableIfNotExistsRe, Callback: rewriteCreateTableIfNotExists},
}
