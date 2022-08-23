package ript

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

//go:generate go run genconfig.go

type Config struct {
	Dryrun       bool
	NoClobber    bool
	Verbose      bool
	LogApis      bool
	VersionToUse int
}

var theConfig *Config

func init() {
	theConfig = &Config{}

	//theConfig.NoClobber = true
	//theConfig.Verbose = true
	//theConfig.LogApis = true
	//theConfig.VersionToUse = 0 // 0 == orig
	theConfig.VersionToUse = 1 // 1 == tarfile
	//theConfig.VersionToUse = 2 // 2 == old, but in walk callback
}

func NewConfig(dryrun bool) *Config {
	return &Config{
		Dryrun: dryrun,
	}
}

func (c *Config) SetSafe(safe bool) {
	c.NoClobber = safe
}

func (c *Config) SetVerbose(verbose bool) {
	c.Verbose = verbose
}

func ConfigDryRun() bool {
	return theConfig.Dryrun
}

func ConfigNoClobber() bool {
	return theConfig.NoClobber
}

func ConfigClobber() bool {
	return !theConfig.NoClobber
}

func ConfigVerbose() bool {
	return theConfig.Verbose
}

func ConfigLogApis() bool {
	return theConfig.LogApis
}

func ConfigUseVersion() int {
	return theConfig.VersionToUse
}

func ConfigIf(n int) bool {
	if n == 0 {
		return false
	}
	return true
}
