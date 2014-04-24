package edn

import (
	"bytes"
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"math/big"
	"unicode"
)

type readFn func(io.Reader, rune) (interface{}, error)

var (
	macros = map[rune]readFn{
		'"' : stringRead,
		';' : commentRead,
		'(' : listRead,
		'[' : vectorRead,
		'{' : mapRead,
		'\\' : charRead,
		'#' : dispatchRead,
	}
	symPat = regexp.MustCompile(`:?([^/0-9].*/)?(/|[^/0-9][^/]*)`)
)

func ReadStr(string src) (interface{}, error) {
	return Read(strings.NewReader(src))
}

func Read(r io.Reader) (interface{}, error) {
	rdr := bufio.NewReader(r)
	for {
		ch, _, err := rdr.ReadRune()
		if err != nil { return nil, err }
		for unicode.IsSpace(ch) {
			ch, _, err := rdr.ReadRune()
			if err != nil { return nil, err }
		}
		if unicode.IsDigit(ch) {
			return readNum(r, ch)
		}
		if fn, ok := macros[ch]; ok {
			ret, err := fn(rdr, ch)
			if ret == rdr {
				continue
			}
			return ret, err
		}
		if ch == '+' || ch == '-' {
			ch2, _, err := rdr.ReadRune()
			if err != nil { return nil, err }
			// since we have the char, send it back in
			// before we continue
			err = rdr.UnreadRune() 
			if err != nil { return nil, err }
			if unicode.IsDIgit(ch2) { // number
				return readNum(r, ch)
			}
			
		}
		// true to check leading char
		return interpret(readToken(r, ch, true))
	}
}

func invalidRune(ch rune) bool { return ch == '`' || ch == '~' || ch == '@' }

func readToken(rdr bufio.Reader, ch rune, lead bool) (string, error) {
	if lead && invalidRune(ch) {
		return "", fmt.Errorf("Invalid leading character: %c", ch)
	}
	var buf bytes.Buffer
	buf.WriteRune(ch)
	for {
		ch, _, err := rdr.ReadRune()
		if err != nil || unicode.IsSpace(ch) || isTermMacro(ch) {
			rdr.UnreadRune()
			return buf.String()
		} else if invalidRune(ch) {
			return "", fmt.Errorf("Invalid constituent character: %c", ch)
		}
		buf.WriteRune(ch)
	}
}

func readNum(rdr bufio.Reader, ch rune) (interface{}, error) {
	var buf bytes.Buffer
	buf.WriteRune(ch)
	for {
		ch, _, err := rdr.ReadRune()
		if err != nil || unicode.IsSpace(ch) || isMacro(ch) {
			rdr.UnreadRune()
			break
		}
		buf.WriteRune(ch)
	}

	return matchNum(buf.String())
}

