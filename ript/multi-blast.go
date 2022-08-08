package ript

/* Copyright © 2022 sparksb -- MIT (see LICENSE file) */
/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

type MultiBlast struct {
	Fs *fs.FS

	dirmask int
	umasked bool
}

func NewMultiBlast(dir string) *MultiBlast {
	dirFS := os.DirFS(dir)

	return &MultiBlast{
		Fs:      &dirFS,
		dirmask: 0775, /* or 755, or 750 */
	}
}

func (m *MultiBlast) cpMany(filenames map[string]string) error {

	if ConfigDryRun() {
		for src, dest := range filenames {
			fmt.Printf("cp '%s' '%s'\n", src, dest)
		}

		return nil
	}

	var totalBytes int64
	for src, dest := range filenames {
		num, err := m.cp(src, dest)
		if err != nil {
			return err
		}
		totalBytes += num
	}

	if ConfigVerbose() {
		fmt.Printf("[%9d bytes total] \n", totalBytes)
	}

	return nil
}

func (m *MultiBlast) cp(src, dest string) (int64, error) {

	destDir := filepath.Dir(dest)
	_, err := m.mkdirpify(destDir)
	if err != nil {
		return 0, err
	}

	// TODO: If exists and user wants warning, do it

	return m.copyBytes(src, dest)
}

func (m *MultiBlast) copyBytes(src, dest string) (int64, error) {
	if ConfigNoClobber() {
		destInfo, err := fs.Stat(*m.Fs, dest)
		//fmt.Println(os.IsExist(err), os.IsNotExist(err), destInfo)
		if err == nil && destInfo != nil && destInfo.Size() > 0 {
			return 0, fmt.Errorf("SafeModeButDestExists %s", dest)
		}
	}

	info, err := fs.Stat(*m.Fs, src)
	if err != nil {
		return 0, err
	}
	srcSize := info.Size()

	destFd, err := os.Create(dest)
	if err != nil {
		log.Panic(err)
	}
	defer destFd.Close()

	srcFd, err := os.Open(src)
	if err != nil {
		log.Panic(err)
	}
	defer srcFd.Close()

	numWritten, err := io.Copy(destFd, srcFd)
	if err != nil {
		log.Panic(err)
	}
	assertMsg(numWritten == srcSize, fmt.Sprintf("cp did not copy the right number of bytes: %d vs %d. (%s -> %s)", numWritten, srcSize, src, dest))

	if ConfigVerbose() {
		fmt.Printf("[%9d bytes] written to %s\n", numWritten, dest)
	}

	return numWritten, nil
}

func (m *MultiBlast) mkdirpify(dir string) (bool, error) {
	dirInfo, err := fs.Stat(*m.Fs, dir)
	//fmt.Println(os.IsExist(err), os.IsNotExist(err))
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	created := false
	if dirInfo == nil && os.IsNotExist(err) {
		m.umaskify()
		err = os.MkdirAll(dir, os.FileMode(m.dirmask))
		if err != nil && !os.IsExist(err) {
			log.Panic(err)
		}
		created = !os.IsExist(err)
	}

	return created, nil
}

func (m *MultiBlast) umaskify() {
	if !m.umasked {
		_ = syscall.Umask(0)
		//fmt.Printf("Old umask: %d\n", old)
		m.umasked = true
	}
}
