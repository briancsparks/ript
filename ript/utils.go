package ript

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func first(m map[string]string) (string, string, bool) {
	for k, v := range m {
		return k, v, true
	}
	return "", "", false
}

func asserter(test bool) bool {
	if !test {
		breakout("", true)
	}
	return !test
}

func assertMsg(test bool, msg string) {
	if !test {
		breakout(msg, false)
	}
}

func assert(test bool) {
	if !test {
		breakout("", false)
	}
}

func breakout(msg string, quiet bool) {
	if !quiet {
		fmt.Printf("  ------------ BREAKOUT!! %v !!\n", msg)
	}
}

func joinMap(m *map[string]string) string {
	result := ""
	for k, v := range *m {
		result += fmt.Sprintf("%-32s: %s\n", k, v)
	}
	return result
}

func trimPrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

func TrimPathPrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		_, rest, found := strings.Cut(s, "/")
		if found {
			return rest
		}
		return s[len(prefix)+1:]
	}
	return s
}

func listHasPrefix(s string, l []string) bool {
	for _, s2 := range l {
		if strings.HasPrefix(s, s2) {
			return true
		}
	}
	return false
}

func Cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	return cwd
}

func mkdirp(destPath string, perm os.FileMode) error {
	if ConfigVerbose() {
		fmt.Printf("mkdir -p : %x :%s\n", perm, destPath)
	}

	err := os.MkdirAll(destPath, perm)
	if err != nil {
		log.Panic(err)
	}

	return nil
}
