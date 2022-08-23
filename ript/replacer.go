package ript

import (
	"bufio"
	"fmt"
	"io"
)

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

type Replacer struct {
	r            io.Reader
	replacements map[string]string

	scanner  *bufio.Scanner
	numLines int
}

func NewReplacer(r io.Reader, replacements map[string]string) *Replacer {
	rplr := Replacer{
		r:            r,
		replacements: replacements,
	}
	rplr.scanner = bufio.NewScanner(rplr.r)

	return &rplr
}

func (rplr *Replacer) Read(p []byte) (int, error) {

	count := 0
	xerr := rplr.scanner.Err()

	var line string
	var linelen int
	if rplr.scanner.Scan() {
		line = rplr.scanner.Text()
		linelen = len(line)
		line = replaceEmAll(line, rplr.replacements)
		count = copy(p, line+"\n")
		rplr.numLines += 1
		fmt.Printf("%04d count: %v %v %s\n", rplr.numLines, linelen, count, line)

	} else {
		// Scan() failed, end
		if err := rplr.scanner.Err(); err != nil {
			return 0, err
		}

		return 0, io.EOF
	}
	xerr = rplr.scanner.Err()

	_ = xerr
	return count, nil
}
