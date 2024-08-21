package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name: "plausible-ddl-proxy",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "addr", Value: "localhost:8000"},
			&cli.StringFlag{Name: "target", Value: "http://localhost:8123"},
			&cli.BoolFlag{Name: "disable-rewrites", Value: false},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run(ctx *cli.Context) error {
	var (
		rctx = rewriteContext{
			ZkPathPrefix: "/clickhouse/{cluster}/tables/{shard}/{database}/",
			Cluster:      "{cluster}",
			Replica:      "{replica}",
		}

		rws []rewrites
	)

	sctx, cancel := signal.NotifyContext(ctx.Context, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if !ctx.Bool("disable-rewrites") {
		rws = []rewrites{
			engineRewrites,
			onClusterRewrites,
		}
	}

	var server = http.Server{
		Addr:    ctx.String("addr"),
		Handler: http.HandlerFunc(makeRewriteHandler(rctx, ctx.String("target"), rws)),
	}

	slog.Info("server running", "addr", ctx.String("addr"))

	go func() {
		<-sctx.Done()
		server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err == http.ErrServerClosed {
		return nil
	} else {
		return err
	}
}

func makeRewriteHandler(rctx rewriteContext, target string, rewrites []rewrites) http.HandlerFunc {
	var client http.Client

	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(target)

		if err != nil {
			fail(w, err)
			return
		}

		u.Path = r.URL.Path
		u.RawQuery = r.URL.Query().Encode()
		r.URL = u
		r.RequestURI = ""
		r.RemoteAddr = ""
		r.Host = r.URL.Host
		r.ContentLength = 0

		body, err := io.ReadAll(r.Body)

		if err != nil {
			fail(w, err)
			return
		}

		body, err = rewriteBody(rctx, rewrites, body)

		if err != nil {
			fail(w, err)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body))
		resp, err := client.Do(r)

		if err != nil {
			fail(w, err)
			return
		}

		for k, vs := range resp.Header {
			if len(vs) == 0 {
				continue
			}

			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(resp.StatusCode)

		var buf = make([]byte, 1024)
		_, err = io.CopyBuffer(w, resp.Body, buf)

		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func fail(w http.ResponseWriter, err error) {
	slog.Error(err.Error())
	w.WriteHeader(500)
}

func rewriteBody(rctx rewriteContext, rewrites []rewrites, body []byte) ([]byte, error) {
	var (
		old = string(body)
		new = old
	)

	if !strings.Contains(strings.ToLower(old), "create table") {
		slog.Info("SKIPPED", "query", old)
		return body, nil
	}

	for _, rws := range rewrites {
		new = rws.Run(rctx, new)
	}

	slog.Info("REWRITTEN", "old", old, "new", new)

	return []byte(new), nil
}
