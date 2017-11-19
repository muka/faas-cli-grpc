package api

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func addDir(source string, archive *zip.Writer) error {
	baseDir := filepath.Base(source)
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return addDir(filepath.Join(source, info.Name()), archive)
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()
		_, err = io.Copy(writer, file)

		return err
	})

	return err
}

//Zip compress a directory to ZIP
func Zip(target string, sources ...string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	for _, source := range sources {

		info, err := os.Stat(source)
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if err := addDir(source, archive); err != nil {
				return err
			}
		} else {

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			file, err := os.Open(source)
			if err != nil {
				file.Close()
				return err
			}
			_, err = io.Copy(writer, file)
			if err != nil {
				file.Close()
				return err
			}

			file.Close()
		}
	}
	return nil
}

//Unzip uncompress a ZIP archive
func Unzip(archive []byte, target string) error {
	buf := bytes.NewReader(archive)
	reader, err := zip.NewReader(buf, buf.Size())
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
