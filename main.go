package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var (
	debug   bool
	version bool

	file string
	user string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [ - | host... ]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.BoolVarP(&version, "version", "v", false, "display version information")
	flag.BoolVarP(&debug, "debug", "d", false, "enable debug logging")

	flag.StringVarP(&file, "file", "f", "shops.yml", "configuration file")
	flag.StringVarP(&user, "user", "u", "root", "default user to connect to")
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

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Error: must supply one or more hosts or `-` to read hosts from stdin")
		flag.Usage()
		os.Exit(1)
	}

	hostaddr := flag.Arg(0)
	client, session, err := connectToHost(user, hostaddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting to host %s: %s", hostaddr, err)
		os.Exit(2)
	}
	defer client.Close()

	out, err := session.CombinedOutput("uptime")
	if err != nil {
		log.WithError(err).Error("error running command")
	} else {
		fmt.Println(string(out))
	}
}
