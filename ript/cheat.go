package ript

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
	"bufio"
	"fmt"
	"github.com/spf13/pflag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func Cheat(srcDir, destDir string) error {

	var nocopy map[string]string
	var keys []string
	var envkeys map[string]string

	riptfilename := filepath.Join(srcDir, "riptfile.yaml")
	nocopy, keys, envkeys, err := readRiptfile(riptfilename)
	if err != nil {
		return err
	}

	nocopy["riptfile.yaml"] = "riptfile.yaml"

	_, _, _ = nocopy, keys, envkeys

	// TODO: Check that all keys and envkeys have associated values           ------
	var flagSet = pflag.NewFlagSet("ript", pflag.ContinueOnError)
	_ = flagSet

	var args map[string]*string = make(map[string]*string)
	for _, key := range keys {
		args[key] = flagSet.String(key, "", "")
	}

	err = flagSet.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	var replacements map[string]string = envkeys
	for key, value := range args {
		if len(*value) == 0 {
			fmt.Printf("Missing arg: --%s=%s\n", key, key)
			return fmt.Errorf("Missing arg: %s\n", key)
		}
		replacements[fmt.Sprintf("ript%s", key)] = *value
	}
	//                                                                        --------------------------

	filenames := map[string]string{}
	dirnames := map[string]string{}

	//fmt.Printf("srcDir: %v\n", srcDir)
	//fmt.Printf("destDir: %v\n", destDir)

	srcFS := os.DirFS(srcDir)
	err = fs.WalkDir(srcFS, ".", func(shortPath string, d fs.DirEntry, err error) error {
		if d == nil {
			return err
		}

		// Save a bunch of work
		if d.IsDir() {
			if shortPath == "." {
				return nil
			}
			if shortPath == ".git" || shortPath == "node_modules" || shortPath == ".idea" {
				return fs.SkipDir
			}
		} else {
			_, present := nocopy[shortPath]
			if present {
				return nil
			}
		}

		// Figure out the destination dirs, and make them - make sure to rename things when necessary
		shortDestPath := shortPath
		for oldStr, newStr := range replacements {
			shortDestPath = strings.ReplaceAll(shortDestPath, oldStr, newStr)
		}

		srcPath := filepath.Clean(filepath.Join(srcDir, shortPath))
		destPath := filepath.Clean(filepath.Join(destDir, shortDestPath))
		//fmt.Printf("srcPath: %v, destPath: %v\n", srcPath, destPath)

		//info, _ := d.Info()
		//fmt.Printf("%s: dirent: %v; err: %v\n", shortPath, info, err)

		if d.IsDir() {
			dirnames[srcPath] = destPath
		} else {
			filenames[srcPath] = destPath
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	//fmt.Printf("dirnames:\n%v\n", joinMap(&dirnames))
	//fmt.Printf("filenames:\n%v\n", joinMap(&filenames))

	if ConfigClobber() {
		// Make dest dirs
		_ = syscall.Umask(0)
		for _, dirname := range dirnames {
			err := os.MkdirAll(dirname, os.FileMode(0755)) // TODO: Magic number
			if err != nil && !os.IsExist(err) {
				log.Panic(err)
			}
		}
	}

	if ConfigClobber() {
		// Copy files
		for src, dest := range filenames {
			func(src, dest string) {
				srcFd, err := os.Open(src)
				if err != nil {
					log.Panic(err)
				}
				defer srcFd.Close()

				//destFd, err := os.Create(dest)
				destFd, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // TODO: Magic number
				if err != nil {
					log.Panic(err)
				}
				defer destFd.Close()

				datawriter := bufio.NewWriter(destFd)

				scanner := bufio.NewScanner(srcFd)
				for scanner.Scan() {
					line := scanner.Text()
					line = replaceEmAll(line, replacements)
					datawriter.WriteString(line + "\n")
				}
				if err := scanner.Err(); err != nil {
					log.Panic(err)
				}
				datawriter.Flush()

			}(src, dest)
		}
	}

	return nil
}

func replaceEmAll(str string, replacements map[string]string) string {
	result := str
	for oldStr, newStr := range replacements {
		result = strings.ReplaceAll(result, oldStr, newStr)
	}

	return result
}
