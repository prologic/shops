package main

import (
	"fmt"
	"os"
	"strings"

	scp "github.com/hnakamur/go-scp"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"
)

var (
	debug   bool
	version bool

	file  string
	user  string
	port  int
	quiet bool
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
	flag.BoolVarP(&quiet, "quiet", "q", false, "quiet operationg (no command output")
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

	config, err := readConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config %s: %s", file, err)
		os.Exit(2)
	}

	for _, host := range flag.Args() {
		hostaddr := parseHost(host, port)
		fmt.Printf("%s:\n", hostaddr)

		client, _, err := connectToHost(user, hostaddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, " error connecting to host %s: %s\n", hostaddr, err)
			continue
		}

		scpClient := scp.NewSCP(client)
		for _, file := range config.Files {
			fileInfo, err := os.Stat(file.Source)
			if err != nil {
				log.WithError(err).Error("error getting file info")
				continue
			}

			if fileInfo.IsDir() {
				err = scpClient.SendDir(file.Source, file.Target, nil)
			} else {
				err = scpClient.SendFile(file.Source, file.Target)
			}

			if err == nil {
				fmt.Printf(" %s ✅\n", file)
			} else {
				fmt.Printf("%s ERR (%s)\n", file, err)
			}
		}

		for _, item := range config.Items {
			out, err := executeCommand(item.Check, hostaddr, client)
			if err == nil {
				if quiet {
					fmt.Printf(" %s ✅\n", item)
				} else {
					fmt.Printf(" %s ✅ -> %s\n", item, strings.TrimSpace(out))
				}
				continue
			}

			if exitError, ok := err.(*ssh.ExitError); ok && exitError.ExitStatus() != 0 {
				out, err := executeCommand(item.Action, hostaddr, client)
				if err == nil {
					if quiet {
						fmt.Printf(" %s ✅\n", item)
					} else {
						fmt.Printf(" %s ✅ -> %s\n", item, strings.TrimSpace(out))
					}
					continue
				}

				if exitError, ok := err.(*ssh.ExitError); ok && exitError.ExitStatus() != 0 {
					fmt.Printf("%s ERR (Status: %d Output: %s)\n", item, exitError.ExitStatus(), out)
				}
			} else {
				log.WithError(err).Errorf("error running check %s against %s", item, hostaddr)
				fmt.Printf("%s ERR (Status: %d Output: %s)\n", item, exitError.ExitStatus(), out)
			}
		}
	}
}
