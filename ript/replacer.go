package ript

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
	"bufio"
	"io"
)

type Replacer struct {
	origReader   io.Reader
	replacements map[string]string

	//reader *bufio.Reader
	scanner  *bufio.Scanner
	numLines int
}

func NewReplacer(r io.Reader, replacements map[string]string) *Replacer {
	rplr := Replacer{
		origReader:   r,
		replacements: replacements,
	}
	rplr.scanner = bufio.NewScanner(rplr.origReader)
	//rplr.reader = bufio.NewReader(rplr.origReader)

	return &rplr
}

func (rplr *Replacer) Read(p []byte) (int, error) {

	count := 0

	xerr := rplr.scanner.Err()

	var line string
	//var linelen int
	if rplr.scanner.Scan() {
		line = rplr.scanner.Text()
		//linelen = len(line)
		line = replaceEmAll(line, rplr.replacements)
		count = copy(p, line+"\n")
		rplr.numLines += 1
		//fmt.Printf("%04d count: %v %v %s\n", rplr.numLines, linelen, count, line)

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

//func (rplr *Replacer) Read(p []byte) (int, error) {
//
//  count := 0
//
//  line, err := rplr.reader.ReadString('\n')
//  if err != nil {
//    if err == io.EOF {
//      return 0, io.EOF
//    }
//
//    log.Panic(err)
//  }
//
//  linelen := len(line)
//  line = replaceEmAll(line, rplr.replacements)
//  count = copy(p, line)
//  rplr.numLines += 1
//  fmt.Printf("%04d count: %v %v %s", rplr.numLines, linelen, count, line)
//
//  return count, nil
//}
