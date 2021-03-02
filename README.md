# shops - Shell Operations

![GitHub All Releases](https://img.shields.io/github/downloads/prologic/shops/total)

![](https://github.com/prologic/shops/workflows/Go/badge.svg)
![](https://github.com/prologic/shops/workflows/ReviewDog/badge.svg)

[![Go Report Card](https://goreportcard.com/badge/prologic/shops)](https://goreportcard.com/report/prologic/shops)
[![codebeat badge](https://codebeat.co/badges/15fba8a5-3044-4f40-936f-9e0f5d5d1fd9)](https://codebeat.co/projects/github-com-prologic-shops-master)
[![GoDoc](https://godoc.org/github.com/prologic/shops?status.svg)](https://godoc.org/github.com/prologic/shops)
[![GitHub license](https://img.shields.io/github/license/prologic/shops.svg)](https://github.com/prologic/shops)

`shops` is a simple command-line tool written in [Go](https://golang.org)
that helps you simplify the way you manage configuration across a set of
machines. `shops` is your configuration management tool of choice when Chef,
Puppet, Ansible are all too complicated and all you really want to do is
run a bunch of regular shell against a set of hosts.

## Getting Started

To install `shops` you can either run `go get` directly:

```#!console
go get git.mills.io/prologic/shops
```

> __NOTE:__ Be sure to have `$GOBIN` in your `$PATH`. See `go env`.

Or grab the source code and build:

```#!console
git clone https://git.mills.io/prologic/shops.git
cd shops
go build
```

And optionally run `go install` to place the binary `shops` in your `$GOBIN`.

## Usage

Using `shops` is quite simple. The basic usage is as follows:

```#!console
shops -f /path/to/config.yml <host1> <host2> <hostN>
```

For example running the included `test.yml` configuration file at the root of
the source code repository here against a typical Linux server:

```#!console
shops -f test.yml 10.0.0.50:22
```

Will perform the will perform the following:

- Copy `README.md` to `/root/README.md` on the server
- Ensure `/tmp/foo` exists
- Check the uptime of the server and display it.

Example:

```#!console
$ ./shops -f test.yml 10.0.0.50:22
10.0.0.50:22:
 Ensure /root/foo exists ✅ -> /root/foo
 Ensure sshbox is running ✅ ->
 Check Uptime ✅ -> 04:01:12 up 19 days,  2:18,  0 users,  load average: 0.00, 0.00, 0.00
```

## Configuration Specification

Currently the configuration specification is a simple YAML file that consists
of a number of top-level keys:

- `version` -- Which for the moment is ignored, but _might_ be used to version
  the configuration file for future enhancements in a backwards compatible way.
- `files` -- Declares one or more files or directories to be copied to each
  target host. Directories are copied recursively. Currently no checks are
  performed, but this is planned.
- `items` -- One or more items of configuration to be applied to each target
  host. Each item declares a "check" and "action". Checks and actions are
  written in regular shell. If a check fails, the action is run to correct the
  failed state. If all checks pass, no actions are run.

## License

`shops` is licensed under the terms of the [MIT License](/LICENSE)
