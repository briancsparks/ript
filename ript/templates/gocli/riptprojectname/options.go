package riptprojectname

/* Copyright © 2022 RIPTENV_IMFO_USER_IRL_NAME <RIPTENV_IMFO_GITHUB_USEREMAIL> -- MIT (see LICENSE file) */

type DebugOptions struct {
	VerboseLevel int
}

var Debug DebugOptions

func init() {
	Debug.VerboseLevel = 0
}
