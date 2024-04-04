package main

import (
	"bufio"
	"embed"
	"net/url"
	"strings"
)

//go:embed tlds-alpha-by-domain.txt
var domainTXT embed.FS // https://data.iana.org/TLD/tlds-alpha-by-domain.txt

func LoadDomainList() (map[string]struct{}, error) {
	file, err := domainTXT.Open("tlds-alpha-by-domain.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var linesMap = make(map[string]struct{})
	var lineNumber int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if lineNumber == 0 {
			lineNumber++
			continue
		}
		linesMap[scanner.Text()] = struct{}{}
		lineNumber++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return linesMap, nil
}

func ParseHost(host string) (string, error) {
	parse, err := url.Parse(host)
	if err != nil {
		return "", err
	}
	h := parse.Host
	split := strings.Split(h, ".")
	if len(split) == 0 {
		return "", nil
	}
	if len(split) == 1 {
		return split[0], nil
	}
	return split[1], nil
}
