package edn

import (
	"bytes"
	"bufio"
	"io"
	"strings"
	"math/big"
	"unicode"
)

type readFn func(io.Reader, rune) interface{}

func Read(rdr bufio.Reader) (interface{}, error) {
	for {
		ch, _, err := rdr.ReadRune()
		if err != nil {
			return nil, err
		}
		if unicode.IsDigit(ch) {
			return readNum(r, ch)
		}
		
			