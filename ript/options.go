package ript

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

type DebugOptions struct {
	VerboseLevel int
}

var Debug DebugOptions

func init() {
	Debug.VerboseLevel = 0
}
