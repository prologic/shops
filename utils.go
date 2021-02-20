package main

import (
	"fmt"
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
