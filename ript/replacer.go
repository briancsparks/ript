package ript

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
	"bufio"
	"io"
)

// -------------------------------------------------------------------------------------------------------------------

type Replacer struct {
	origReader   io.Reader
	replacements map[string]string

	scanner  *bufio.Scanner
	numLines int
}

// -------------------------------------------------------------------------------------------------------------------

func NewReplacer(r io.Reader, replacements map[string]string) *Replacer {
	rplr := Replacer{
		origReader:   r,
		replacements: replacements,
	}
	rplr.scanner = bufio.NewScanner(rplr.origReader)

	return &rplr
}

// -------------------------------------------------------------------------------------------------------------------

func (rplr *Replacer) Read(p []byte) (int, error) {

	count := 0

	var line string
	if rplr.scanner.Scan() {
		line = rplr.scanner.Text()
		line = replaceEmAll(line, rplr.replacements)
		count = copy(p, line+"\n")
		rplr.numLines += 1

	} else {
		// Scan() failed, end
		if err := rplr.scanner.Err(); err != nil {
			return 0, err
		}

		return 0, io.EOF
	}

	return count, nil
}
