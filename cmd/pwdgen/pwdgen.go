package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"io"
	"log"
	"math/big"

	"github.com/atotto/clipboard"
)

var (
	length = flag.Int("length", 16, "生成的密码长度，【6, 2048】")
	level  = flag.Int("level", 4, "生成的密码强度等级，数字越大，强度越高，【1, 4】")
)

func main() {

	flag.Parse()

	s, err := FullPassword(*level, *length)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	if err := clipboard.WriteAll(s); err != nil {
		log.Fatalf("err: %v\n", err)
	}
}

const (
	// LowerLetters is the list of lowercase letters.
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"

	// UpperLetters is the list of uppercase letters.
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Digits is the list of permitted digits.
	Digits = "0123456789"

	// Symbols is the list of symbols.
	Symbols = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
)

const (
	_defaultLength          = 16
	_defaultNumLowerLetters = 4
	_defaultNumUpperLetters = 4
	_defaultNumDigits       = 4
	_defaultNumSymbols      = 4
)

type FullPasswordConf struct {
	Length          int
	NumLowerLetters int
	NumUpperLetters int
	NumDigits       int
	NumSymbols      int
}

func SetLevel(level, length int) FullPasswordConf {

	var fullConf FullPasswordConf

	if level == 1 {
		fullConf.NumDigits = length
	} else if level == 2 {
		fullConf.NumLowerLetters = length / 2
		fullConf.NumDigits = length - fullConf.NumLowerLetters
	} else if level == 3 {
		fullConf.NumDigits = length / 3
		fullConf.NumUpperLetters = (length - fullConf.NumDigits) / 2
		fullConf.NumLowerLetters = length - fullConf.NumDigits - fullConf.NumUpperLetters
	} else if level == 4 {
		fullConf.NumDigits = length / 5
		fullConf.NumUpperLetters = length / 4
		fullConf.NumLowerLetters = length / 4
		fullConf.NumSymbols = length - fullConf.NumDigits - fullConf.NumUpperLetters - fullConf.NumLowerLetters
	} else {
		fullConf.Length = _defaultLength
		fullConf.NumLowerLetters = _defaultNumLowerLetters
		fullConf.NumUpperLetters = _defaultNumUpperLetters
		fullConf.NumDigits = _defaultNumDigits
		fullConf.NumSymbols = _defaultNumSymbols
	}
	return fullConf
}

func FullPassword(level, length int) (string, error) {

	if level < 1 || level > 4 {
		return "", errors.New("level must range 1-4")
	}

	if length < 6 {
		return "", errors.New("length must >= 6")
	} else if length > 2048 {
		return "", errors.New("length too long")
	}

	var (
		result string
		read   = rand.Reader
	)

	var fullConf = SetLevel(level, length)

	// Characters
	for i := 0; i < fullConf.NumLowerLetters; i++ {
		ch, err := randomElement(read, LowerLetters)
		if err != nil {
			return "", err
		}

		result, err = randomInsert(read, result, ch)
		if err != nil {
			return "", err
		}
	}

	for i := 0; i < fullConf.NumUpperLetters; i++ {
		ch, err := randomElement(read, UpperLetters)
		if err != nil {
			return "", err
		}

		result, err = randomInsert(read, result, ch)
		if err != nil {
			return "", err
		}
	}

	// Digits
	for i := 0; i < fullConf.NumDigits; i++ {
		d, err := randomElement(read, Digits)
		if err != nil {
			return "", err
		}

		result, err = randomInsert(read, result, d)
		if err != nil {
			return "", err
		}
	}

	// Symbols
	for i := 0; i < fullConf.NumSymbols; i++ {
		sym, err := randomElement(read, Symbols)
		if err != nil {
			return "", err
		}

		result, err = randomInsert(read, result, sym)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

// randomInsert randomly inserts the given value into the given string.
func randomInsert(reader io.Reader, s, val string) (string, error) {
	if s == "" {
		return val, nil
	}

	n, err := rand.Int(reader, big.NewInt(int64(len(s)+1)))
	if err != nil {
		return "", err
	}
	i := n.Int64()
	return s[0:i] + val + s[i:], nil
}

// randomElement extracts a random element from the given string.
func randomElement(reader io.Reader, s string) (string, error) {
	n, err := rand.Int(reader, big.NewInt(int64(len(s))))
	if err != nil {
		return "", err
	}
	return string(s[n.Int64()]), nil
}
