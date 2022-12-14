package ript

import (
	"github.com/goccy/go-yaml"
	"os"
	"strings"
)

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

func readRiptfile(riptfilename string) (map[string]string, []string, map[string]string, []string, error) { /* nocopy, keys, envkeys, missingenv */
	nocopyPath, err := yaml.PathString("$.nocopy")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	keysPath, err := yaml.PathString("$.keys")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	envkeysPath, err := yaml.PathString("$.envkeys")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	riptfileData, err := os.ReadFile(riptfilename)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var nocopy0 []string
	if err = nocopyPath.Read(strings.NewReader(string(riptfileData)), &nocopy0); err != nil {
		return nil, nil, nil, nil, err
	}
	var nocopy map[string]string = make(map[string]string)
	for _, s := range nocopy0 {
		nocopy[s] = s
	}

	var keys []string
	if err = keysPath.Read(strings.NewReader(string(riptfileData)), &keys); err != nil {
		return nocopy, nil, nil, nil, err
	}

	var envkeys map[string]string
	if err = envkeysPath.Read(strings.NewReader(string(riptfileData)), &envkeys); err != nil {
		return nocopy, keys, nil, nil, err
	}

	var missingenv []string
	for name, defValue := range envkeys {
		value, present := getEnv(name, defValue)
		envkeys[name] = value

		if !present {
			missingenv = append(missingenv, name)
		}
	}

	return nocopy, keys, envkeys, missingenv, err
}

func getEnv(key, defValue string) (string, bool) { /* value, isPresent */
	value, present := os.LookupEnv(key)
	if present {
		return value, true
	}

	if strings.HasPrefix(key, "RIPTENV_") {
		key2 := key[8:]
		value, present = os.LookupEnv(key2)
		if present {
			return value, true
		}

		key2 = "RIPT_" + key2
		value, present = os.LookupEnv(key2)
		if present {
			return value, true
		}
	}

	return defValue, false
}
