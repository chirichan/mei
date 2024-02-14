package main

import (
	"bufio"
	"embed"
	"io"
	"net/url"
	"strings"
	"unicode"

	"github.com/gocarina/gocsv"
	"golang.org/x/net/idna"
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

func IsASCII(s string) bool {
	for _, v := range s {
		if !unicode.Is(unicode.ASCII_Hex_Digit, v) {
			return false
		}
	}
	return true
}

func EncodeToPunycode(chineseDomain string) (string, error) {
	punycodeDomain, err := idna.ToASCII(chineseDomain)
	if err != nil {
		return "", err
	}
	return strings.ToUpper(punycodeDomain), nil
}

// XyKey version 8
type XyKey struct {
	Version int   `json:"version"`
	Key     []Key `json:"key"`
}

type Key struct {
	Name      string  `json:"name"`
	Account   string  `json:"account"`
	Password  string  `json:"password"`
	Password2 string  `json:"password2"`
	Url       string  `json:"url"`
	Note      string  `json:"note"`
	Extra     []Extra `json:"extra"`
}

type Extra struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// ChromeCSV chrome, edge csv password
type ChromeCSV struct {
	Name     string `json:"name" csv:"name"`
	URL      string `json:"url" csv:"url"`
	Username string `json:"username" csv:"username"`
	Password string `json:"password" csv:"password"`
}

func (c *ChromeCSV) GuessWhat() string {
	return c.Name
}

func LoadChromeCSV(r io.Reader) ([]ChromeCSV, error) {
	var temp []*ChromeCSV
	if err := gocsv.Unmarshal(r, &temp); err != nil {
		return nil, err
	}
	var res []ChromeCSV
	for _, v := range temp {
		res = append(res, *v)
	}
	return res, nil
}

func CSVToXykey(version int, csv []ChromeCSV) XyKey {
	var xyKey = XyKey{
		Version: version,
		Key:     make([]Key, len(csv)),
	}
	for _, v := range csv {
		var key Key
		key.Name = v.Name
		key.Account = v.Username
		key.Password = v.Password
		key.Url = v.URL
		xyKey.Key = append(xyKey.Key, key)
	}
	return xyKey
}
