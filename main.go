package main

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var (
	debug   bool
	version bool

	cont bool
	file string
	user string
	port int
	env  []string
)

const helpText = `
shops runs a spec against one or more targets, targets can be provided as
arguments or read from standard input by supplying a single argument "-".

The syntax targets are:

<type>://[<user>@]<hostname>[:<port>]

For local targets, only the type is required. e.g: local://

For remote (ssh://) targets, if either the user or port is not provided as
part of the target, it defaults to the -u/--user and -p/--port options.

Valid options:
`

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [ TARGET [ TARGET ] ... ]\n", os.Args[0])
		fmt.Fprint(os.Stderr, helpText)
		flag.PrintDefaults()
	}

	flag.BoolVarP(&version, "version", "v", false, "display version information and exit")
	flag.BoolVarP(&debug, "debug", "d", false, "enable debug logging")

	flag.StringVarP(&file, "file", "f", "shops.yml", "spec file")
	flag.StringVarP(&user, "user", "u", "root", "default user for temotee targets")
	flag.IntVarP(&port, "port", "p", 22, "default port for remote targets")
	flag.BoolVarP(&cont, "continue-on-error", "c", false, "continue on errors")
	flag.StringSliceVarP(&env, "env", "e", []string{}, "set environment variables")
}

func main() {
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if version {
		fmt.Printf("shops version %s", FullVersion())
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Error: must supply one or more hosts or `-` to read hosts from stdin")
		flag.Usage()
		os.Exit(1)
	}

	config, err := readConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config %s: %s", file, err)
		os.Exit(2)
	}
	config.SetEnvVars(env)

	var uris []URI

	if flag.NArg() == 1 && flag.Arg(0) == "-" {
		lines, err := readLines(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading hosts: %s", err)
			os.Exit(2)
		}
		uris = ParseURIs(lines, user, strconv.Itoa(port))
	} else {
		uris = ParseURIs(flag.Args(), user, strconv.Itoa(port))
	}

	runner, err := NewGroupRunner(
		uris, config,
		WithContinueOnError(cont),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating runner: %s", err)
		os.Exit(2)
	}

	if err := runner.Run(); err != nil {
		Poo()
		os.Exit(3)
	}

	Pony()
	os.Exit(0)
}
