//go:build ignore

package main

import (
	"archive/tar"
	"fmt"
	"github.com/briancsparks/ript/ript"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	err := doOneTar("templates", "gocli")
	if err != nil {
		log.Fatal(err)
	}
}

func doOneTar(srcDir, root string) error {
	tarballPath := filepath.Join(srcDir, root+".tar")
	fmt.Printf("tarballPath: %v\n", tarballPath)

	file, err := os.Create(tarballPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := tar.NewWriter(file)
	defer writer.Close()

	srcFS := os.DirFS(srcDir)
	err = fs.WalkDir(srcFS, root, func(shortPath string, dirent fs.DirEntry, err error) error {

		if dirent == nil {
			return err
		}

		if !dirent.IsDir() {
			filename := filepath.Join(srcDir, shortPath)
			pathInTarFile := ript.TrimPathPrefix(shortPath, root)
			err = addFile(filename, pathInTarFile, writer)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func addFile(filePath, pathInTarFile string, writer *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	fmt.Printf("%9d %+v %+v %s\n", stat.Size(), stat.Mode(), stat.ModTime(), pathInTarFile)
	header := &tar.Header{
		Name:    pathInTarFile,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	err = writer.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
}
