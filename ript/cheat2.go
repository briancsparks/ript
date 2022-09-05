package ript

import (
	"archive/tar"
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

// -------------------------------------------------------------------------------------------------------------------
// TODO: Have the cheat command call Cheat2

func Cheat2(tname, destDir string) error {
	if ConfigLogApis() {
		fmt.Printf("API(Cheat): tname: %v, dest: %s\n", tname, destDir)
	}

	//srcDir := "/home/sparksb/go/src/bcs/tryouts/__go-project-template/one"
	srcDir := filepath.Join(MyDirname, "templates", tname)
	return Cheat2B(srcDir, destDir, tname)
}

// -------------------------------------------------------------------------------------------------------------------

func Cheat2B(srcDir, destDir, tname string) error {
	if ConfigLogApis() {
		fmt.Printf("API(CheatB): srcDir: %v, destDir: %s\n", srcDir, destDir)
	}

	riptfilename := filepath.Join(srcDir, "riptfile.yaml")
	nocopy, keys, envkeys, missingenv, err := readRiptfile(riptfilename)
	if err != nil {
		fmt.Printf("Cannot open riptfilename: %v\n  Is your template name (%s) right?\n\n", riptfilename, tname)
		return err
	}

	if len(missingenv) > 0 {
		fmt.Printf("Missing ENVs: %v\n", missingenv)
		return fmt.Errorf("Missing ENVs: %v\n", missingenv)
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
	args["dest"] = flagSet.String("dest", "", "")

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

	if IsProd() {
		return cheatTarfile2(srcDir, destDir, tname, nocopy, &replacements)
	}

	// ================================================
	// Active Development :)

	if ConfigVerbose() {
		fmt.Printf("srcDir: %v\n", srcDir)
		fmt.Printf("destDir: %v\n", destDir)
	}

	err = fs.WalkDir(os.DirFS(srcDir), ".", func(shortPath string, dirent fs.DirEntry, err error) error {
		if ConfigVerbose() {
			fmt.Printf("shortPath: %v\n", shortPath)
		}

		if dirent == nil {
			return err
		}

		fi, err := dirent.Info()
		if err != nil {
			return err
		}

		// Figure out the source and destination dirs
		srcPath := filepath.Clean(filepath.Join(srcDir, shortPath))
		destPath := filepath.Join(destDir, shortPath)
		// Put callback here - we have all info, let callback do any dest path translation

		destPath = filepath.Clean(replaceEmAll2(destPath, &replacements))

		// Save a bunch of work
		if dirent.IsDir() {
			if shortPath == "." {
				mkdirp(destDir, fi.Mode()) // The 'root' destination dir
				return nil
			}

			mkdirp(destPath, fi.Mode())

		} else {
			_, present := nocopy[shortPath]
			if present {
				return nil
			}

			func() {
				srcFd, err := os.Open(srcPath)
				if err != nil {
					log.Panic(err)
				}
				defer srcFd.Close()

				putFile2(srcFd, destPath, fi.Mode(), &replacements)
			}()
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

func cheatTarfile2(srcDir, destDir, tname string, nocopy map[string]string, replacements *map[string]string) error {

	tarfile, err := templates.Open(filepath.Join("templates", tname+".tar"))
	if err != nil {
		log.Panic(err)
	}
	defer tarfile.Close()

	mkdirp(destDir, 0755) // TODO: Magic number

	tr := tar.NewReader(tarfile)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			log.Panic(err)
		}

		_, present := nocopy[filepath.Base(hdr.Name)]
		if present {
			continue
		}

		destPath := filepath.Join(destDir, hdr.Name)
		// Put callback here - we have all info, let callback do any dest path translation

		destPath = replaceEmAll2(destPath, replacements)
		mkdirp(filepath.Dir(destPath), 0755) // TODO: Magic Number

		// Copy files
		putFile2(tr, destPath, os.FileMode(hdr.Mode), replacements)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------
// Pass nil for replacements for binary file

func putFile2(srcFd io.Reader, destName string, perm os.FileMode, replacements *map[string]string) {
	var count int64

	if ConfigVerbose() {
		fmt.Printf("putFile2 : %v :%s\n", perm, destName)
	}

	if ConfigClobber() {
		destFd, err := os.OpenFile(destName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
		if err != nil {
			fmt.Printf("While trying to create %s\n", destName)
			log.Panic(err)
		}
		defer destFd.Close()

		if replacements != nil {
			replacer := NewReplacer(srcFd, *replacements)
			count, err = io.Copy(destFd, replacer)

		} else {
			count, err = io.Copy(destFd, srcFd)
		}

		if err != nil {
			log.Panic(err)
		}
		_ = count

	}
}

func replaceEmAll2(str string, replacements *map[string]string) string {
	result := str
	for oldStr, newStr := range *replacements {
		result = strings.ReplaceAll(result, oldStr, newStr)
	}

	return result
}
