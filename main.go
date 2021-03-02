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
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [ - | host... ]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.BoolVarP(&version, "version", "v", false, "display version information")
	flag.BoolVarP(&debug, "debug", "d", false, "enable debug logging")

	flag.StringVarP(&file, "file", "f", "shops.yml", "configuration file")
	flag.StringVarP(&user, "user", "u", "root", "default user to authenticate as")
	flag.IntVarP(&port, "port", "p", 22, "default port to connect to remote host")
	flag.BoolVarP(&cont, "continue-on-error", "c", false, "continue on errors")
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

	runner, err := NewGroupRunner(uris, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating runner: %s", err)
		os.Exit(2)
	}

	runner.Run()

	Pony()
}
