# shops - SHell OPerationS

![GitHub All Releases](https://img.shields.io/github/downloads/prologic/shops/total)

![](https://github.com/prologic/shops/workflows/Go/badge.svg)
![](https://github.com/prologic/shops/workflows/ReviewDog/badge.svg)

[![Go Report Card](https://goreportcard.com/badge/prologic/shops)](https://goreportcard.com/report/prologic/shops)
[![codebeat badge](https://codebeat.co/badges/15fba8a5-3044-4f40-936f-9e0f5d5d1fd9)](https://codebeat.co/projects/github-com-prologic-shops-master)
[![GoDoc](https://godoc.org/github.com/prologic/shops?status.svg)](https://godoc.org/github.com/prologic/shops)
[![GitHub license](https://img.shields.io/github/license/prologic/shops.svg)](https://github.com/prologic/shops)

`shops` is a simple command-line tool written in [Go](https://golang.org)
that helps you simplify the way you manage configuration across a set of
machines.

> `shops` is your configuration management tool of choice when Chef,
> Puppet, Ansible are all too complicated and all you really want to do is
> run a bunch of regular shell against a set of hosts.

`shops` basically lets you (_the oeprator_) run a specification against one
or mote targets. Targets can either be local or remote and the syntax of a
target is:

```
type://[<user>@]<hostname>[:<port>]
```

For local targets, only thee type is required, e.g: `local://`.

For remote targets, the user and port are optional and if not specified in the
target they default to the `-u/--user` and `-p/--port` flags respectively.

## Getting Started

### Install from releases

You can install `shops` by simply downloading the latest version from our
[Release](https://github.com/prologic/shops/releases) page for your platform
and placing the binary in your `$PATH`.

For convenience you can run one of the following which will download and
install  the latest release binary into `/usr/local/bin`:

For Linux x86_64:

```console
curl -s https://api.github.com/repos/prologic/shops/releases/latest | grep browser_download_url | grep Linux_x86_64 | cut -d '"' -f 4 | wget -q -O - -i - | tar -xv shops && mv shops /usr/local/bin/shops
```

For MacOS x86_64:

```console
curl -s https://api.github.com/repos/prologic/shops/releases/latest | grep browser_download_url | grep Darwin_x86_64 | cut -d '"' -f 4 | wget -q -O - -i - | tar -xv shops && mv shops /usr/local/bin/shops
```

### Install from source

To install `shops` you can run `go get` directly:

```#!console
go get github.com/prologic/shops
```

> __NOTE:__ Be sure to have `$GOBIN` (_if not empty_) or your `$GOPATH/bin`
>           in your `$PATH`.
>           See [Compile and install packages and dependencies](https://golang.org/cmd/go/#hdr-Compile_and_install_packages_and_dependencies)

Or grab the source code and build:

```#!console
git clone https://github.com/prologic/shops.git
cd shops
go build
```

And optionally run `go install` to place the binary `shops` in your `$GOBIN`
or `$GOPATH/bin` (_again see note above_).

## Usage

Using `shops` is quite simple. The basic usage is as follows:

```#!console
shops -f /path/to/spec.yml <host1> <host2> <hostN>
```

For example running the included `test.yml` specification file which can be
found in the `./testdata` directory in the source tree as well as other examples:

```#!console
shops -f ./testdata/sample.yml 10.0.0.50
```

Will perform the will perform the following:

- Copy `README.md` to `/root/README.md` on the server
- Ensure `/tmp/foo` exists
- Check the uptime of the server and display it.

Example:

```#!console
$ ./shops -f ./testdata/sample.yml 10.0.0.50
10.0.0.50:22:
  README.md -> /root/README.md ✅
  Ensure /root/foo exists ✅ -> /root/foo
  Ensure sshbox is running ✅ ->
  Check Uptime ✅ -> 13:58:27 up 3 days,  1:38,  0 users,  load average: 0.00, 0.00, 0.00


           ,--,
     _ ___/ /\|
 ,;'( )__, )  ~
//  //   '--;
'   \     | ^
     ^    ^
```

> Yes, it really does print a Pony on success! 🤣

## Specification File Format

The specification file format is a simple YAML file with the following
structure:

```#!yaml
---
veresion: 1

files:
  - source: foo
    target: /tmp/foo

funcs:
  foo: |
    echo "I am a function!"

items:
  - name: Check #1
    check: true
    actino: foo
```

A valid spec consists of a number of top-level keys:

- `version` -- Which for the moment is ignored, but _might_ be used to version
               the configuration file for future enhancements in a backwards
               compatible way.
- `files`   -- Declares one or more files or directories to be copied to each
               target host. Directories are copied recursively. Currently no
               checks are performed, but this is planned.
- `funcs`   -- A mapping of function name to the function's body. There is no
               need to write the function name with `foo() { ... }` as the
               `shops` runners will do this automatically when injecting
               functions into the session(s).
- `items`   -- One or more items of configuration to be applied to each target
               host. Each item declares a "check" and "action". Checks and
               actions are written in regular shell. If a check fails, the
               action is run to correct the failed state. If all checks pass,
               no actions are run.

## License

`shops` is licensed under the terms of the [MIT License](/LICENSE)
