package ript

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
	"fmt"
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
