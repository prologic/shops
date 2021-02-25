package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

func parseHost(host string, port int) string {
	var hostaddr string

	if ok, err := regexp.MatchString(`.*\:[0-9]+`, host); !ok || err != nil {
		hostaddr = fmt.Sprintf("%s:%d", host, port)
	} else {
		hostaddr = host
	}

	return hostaddr
}

func parseHosts(hosts []string, port int) []string {
	var addrs []string

	for _, host := range hosts {
		addrs = append(addrs, parseHost(host, port))
	}

	return addrs
}

func readLines(r io.Reader) (lines []string, err error) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err = scanner.Err()

	return
}
