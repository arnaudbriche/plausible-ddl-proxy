# Plausible-ddl-proxy

## Why ?

## Usage

[Plausible Analytics](https://github.com/plausible/analytics) is an awesome open source, and privacy-friendly web analytics product that relies on [ClickHouse](https://github.com/ClickHouse/ClickHouse) to provide a snappy experience.

While the Community Edition of the product allows for an easy self-hosted setup, it does not offers the option to leverage a replicated ClickHouse setup, potentially limiting read scalability and availability.

The role of this small tool is to intercept DDL queries sent by the Plausible migrate command to ClickHouse and rewrite the table creation queries to make use of Replicated* tables.

Point the proxy to your ClickHouse, and the Plausible migrate command to the proxy.

It has only be tested on v2.1.0 and v2.1.1.

### Build

```sh 
make
```

### Run

```sh
./bin/plausible-ddl-proxy -h
NAME:
   plausible-ddl-proxy - A new cli application

USAGE:
   plausible-ddl-proxy [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --addr value        (default: "localhost:8000")
   --target value      (default: "http://localhost:8123")
   --disable-rewrites  (default: false)
   --help, -h          show help
```