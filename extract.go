package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func extract(filename, ver string) {
	ShowNotification(fmt.Sprintf(strExtrV, ver))
	SetNotificationToolTip(fmt.Sprintf(strExtrV, ver))
	log.Println("Extracting...")
	ext := path.Ext(filename)
	if ext == ".zip" {
		// if zip archive
		unzip(filename, filepath.Join(pgsqlBaseDir, ver))
	} else if ext == (".bz2") || ext == (".gz") {
		// if gz or bz2 archive
		unGzip(filename, filepath.Join(pgsqlBaseDir, ver))
	} else if ext == ".7z" {
		log.Println("7z archive not allowed")
	}
	ShowNotification(fmt.Sprintf(strExtrVF, ver))
	SetNotificationToolTip("")
}

func unzip(filename, dest string) {
	var fpath string
	if filename == "" {
		log.Println("Can't unzip ", filename)
		os.Exit(1)
	}

	reader, err := zip.OpenReader(filename)
	if err != nil {
		log.Printf("Extract error::OpenArchive - %s\n", err)
	}
	defer reader.Close()

	for _, f := range reader.Reader.File {
		zipped, err := f.Open()
		checkErr("Extract error", err)

		fpath = filepath.Join(dest, strings.Replace(f.Name, "pgsql", "", 1))
		// log.Println(fpath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			writer, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, f.Mode())
			if err != nil {
				log.Printf("Extract error::OpenFileFromArchive - %s\n", err)
			}

			if _, err = io.Copy(writer, zipped); err != nil {
				log.Println(err)
				os.Exit(1)
			}
			defer writer.Close()
		}
		defer zipped.Close()
	}
	log.Println("Extracting finished!")
}

func unGzip(sourcefile, dest string) {
	reader, err := os.Open(sourcefile)
	checkErr("In unGzipBzip2 - Open", err)
	defer reader.Close()

	var tarReader *tar.Reader
	var fileReader io.ReadCloser = reader

	if strings.HasSuffix(sourcefile, ".gz") ||
		strings.HasSuffix(sourcefile, ".tgz") {
		gzipReader, err := gzip.NewReader(reader)
		checkErr("In unGzip - NewReader", err)
		tarReader = tar.NewReader(gzipReader)
		defer fileReader.Close()
	} else {
		tarReader = tar.NewReader(reader)
	}

	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			checkErr("Extract error::ReadTarArchive", err)
		}

		filename := ""
		filename = filepath.Join(dest, strings.Replace(header.Name, "pgsql", "", 1))

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(filename, os.FileMode(header.Mode))
			checkErr("In unGzip - MkdirAll", err)
		case tar.TypeReg, tar.TypeRegA:
			writer, err := os.Create(filename)
			checkErr("In unGzip - Create", err)
			io.Copy(writer, tarReader)
			err = os.Chmod(filename, os.FileMode(header.Mode))
			checkErr("In unGzip - Chmod", err)
			writer.Close()
		case tar.TypeSymlink:
			err = os.Symlink(header.Linkname, filename)
		default:
			log.Printf("Unable to untar type: %c in file %s\n", header.Typeflag, filename)
		}
	}
	log.Println("Extracting finished!")
}
