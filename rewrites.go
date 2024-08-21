package main

import (
	"fmt"
	"regexp"
)

type rewriteContext struct {
	Cluster      string
	ZkPathPrefix string
	Replica      string
	TableName    string
}

type rewrite struct {
	Regex    *regexp.Regexp
	Callback func(rewriteContext) func(string) string
}

type rewrites []rewrite

func (rs rewrites) Run(rctx rewriteContext, s string) string {
	rctx.TableName = tableName(s)

	for _, r := range rs {
		if !r.Regex.MatchString(s) {
			continue
		}

		return r.Regex.ReplaceAllStringFunc(s, r.Callback(rctx))
	}

	return s
}

var tableNameRe = regexp.MustCompile(`TABLE\s+(?:IF\s+NOT\s+EXISTS)?\s*([\"\w]+)`)

func tableName(query string) string {
	var m = tableNameRe.FindStringSubmatch(query)

	if len(m) != 2 {
		panic(fmt.Sprintf("cannot find table name in %s (%s)", query, m))
	}

	switch m[1] {
	case "TABLE", "IF", "NOT", "EXISTS":
		panic(fmt.Sprintf("cannot find table name in %s (%s)", query, m))
	default:
		return m[1]
	}
}
