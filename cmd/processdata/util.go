package main

import (
	"io"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readerToString(r io.Reader) string {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	check(err)
	return buf.String()
}
