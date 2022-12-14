package ript

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
	"archive/tar"
	"bufio"
	"embed"
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed all:templates
var templates embed.FS

// -------------------------------------------------------------------------------------------------------------------

func Cheat(tname, destDir string) error {
	if ConfigLogApis() {
		fmt.Printf("API(Cheat): tname: %v, dest: %s\n", tname, destDir)
	}

	//srcDir := "/home/sparksb/go/src/bcs/tryouts/__go-project-template/one"
	srcDir := filepath.Join(MyDirname, "templates", tname)
	return CheatB(srcDir, destDir, tname)
}

// -------------------------------------------------------------------------------------------------------------------

func CheatB(srcDir, destDir, tname string) error {
	if ConfigLogApis() {
		fmt.Printf("API(CheatB): srcDir: %v, destDir: %s\n", srcDir, destDir)
	}

	var nocopy map[string]string
	var keys []string
	var envkeys map[string]string

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
	//                                                                        --------------------------

	if IsProd() || ConfigUseVersion() == 1 {
		return cheatTarfile(srcDir, destDir, tname, nocopy, replacements)
	}

	// ================================================
	// Active Development :)

	filenames := map[string]string{}
	dirnames := map[string]string{}

	srcFS := os.DirFS(srcDir)
	root := "."

	if ConfigVerbose() {
		fmt.Printf("srcDir: %v\n", srcDir)
		fmt.Printf("destDir: %v\n", destDir)
		fmt.Printf("root: %v\n", root)
	}

	err = fs.WalkDir(srcFS, root, func(shortPath string, d fs.DirEntry, err error) error {
		if ConfigVerbose() {
			fmt.Printf("shortPath: %v\n", shortPath)
		}

		if d == nil {
			return err
		}

		fi, err := d.Info()
		_, _ = fi, err

		// Save a bunch of work
		if d.IsDir() {
			if shortPath == "." {
				mkdirp(destDir, fi.Mode())
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

		if ConfigUseVersion() == 2 {
			if d.IsDir() {
				mkdirp(destPath, fi.Mode())

			} else {
				srcFd, err := os.Open(srcPath)
				if err != nil {
					log.Panic(err)
				}

				putFile(srcFd, destPath, replacements, fi.Mode())
			}

		} else {
			if d.IsDir() {
				dirnames[srcPath] = destPath
			} else {
				filenames[srcPath] = destPath
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	if ConfigVerbose() {
		fmt.Printf("dirnames:\n%v\n", joinMap(&dirnames))
		fmt.Printf("filenames:\n%v\n", joinMap(&filenames))
	}

	if ConfigUseVersion() == 0 {
		if ConfigClobber() {
			// Make dest dirs
			//_ = syscall.Umask(0)          /* TODO: Figure out how to put this back */
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
	}

	return nil
}

func putFile(srcFd io.ReadCloser, destName string, replacements map[string]string, perm os.FileMode) {
	defer srcFd.Close()

	if ConfigVerbose() {
		fmt.Printf("putFile : %v :%s\n", perm, destName)
	}

	if ConfigClobber() {
		destFd, err := os.OpenFile(destName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
		if err != nil {
			fmt.Printf("While trying to create %s\n", destName)
			log.Panic(err)
		}
		defer destFd.Close()

		replacer := NewReplacer(srcFd, replacements)
		count, err := io.Copy(destFd, replacer)
		if err != nil {
			log.Panic(err)
		}
		_ = count

	}
}

func cheatTarfile(srcDir, destDir, tname string, nocopy map[string]string, replacements map[string]string) error {
	_, _, _, _, _ = srcDir, destDir, tname, nocopy, replacements

	srcFS := templates

  // embedded FS is unix-path
  tarFilename := "templates/" + tname + ".tar"

	tarfile, err := srcFS.Open(tarFilename)
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

		pathname := filepath.Join(destDir, hdr.Name)
		pathname = replaceEmAll(pathname, replacements)

		if isDir(hdr) {
			err = os.MkdirAll(pathname, os.FileMode(hdr.Mode))
			if err != nil {
				log.Panic(err)
			}
			continue
		}

		_, present := nocopy[filepath.Base(hdr.Name)]
		if present {
			continue
		}

		if ConfigClobber() {
			theDir := filepath.Dir(pathname)
			err = os.MkdirAll(theDir, 0755) // TODO: Magic Number
			if err != nil {
				log.Print(err)
			}

			// Copy files
			func(dest string) {

				destFd, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(hdr.Mode))
				if err != nil {
					log.Panic(err)
				}
				defer destFd.Close()

				replacer := NewReplacer(tr, replacements)
				count, err := io.Copy(destFd, replacer)
				if err != nil {
					log.Panic(err)
				}
				_ = count

			}(pathname)
		}
	}

	return nil
}

func isDir(hdr *tar.Header) bool {
	if strings.HasSuffix(hdr.Name, "/") {
		return true
	}
	return false
}

func replaceEmAll(str string, replacements map[string]string) string {
	result := str
	for oldStr, newStr := range replacements {
		result = strings.ReplaceAll(result, oldStr, newStr)
	}

	return result
}
