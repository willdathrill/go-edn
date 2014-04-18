package edn

import (
	"bufio"
	"io"
	"strings"
	"math/big"
)

type readFn func(io.Reader, rune) interface{}

