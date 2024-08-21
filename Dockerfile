FROM golang:1.22-bullseye as build

WORKDIR /work

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod go mod download -x
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build make

FROM gcr.io/distroless/base-debian11

COPY --from=build /work/bin/* /

ENTRYPOINT ["/plausible-ddl-proxy"]

