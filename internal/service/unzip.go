// sourced from http://blog.ralch.com/articles/golang-working-with-zip/
//TODO is the license on this code permissable enough for me to use it here?
// seems generic enough to be usable
package service

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// func unzip(archive, target string) error {
// 	reader, err := zip.OpenReader(archive)
// 	if err != nil {
// 		return err
// 	}

// 	if err := os.MkdirAll(target, 0755); err != nil {
// 		return err
// 	}

// 	for _, file := range reader.File {
// 		path := filepath.Join(target, file.Name)
// 		if file.FileInfo().IsDir() {
// 			os.MkdirAll(path, file.Mode())
// 			continue
// 		}

// 		fileReader, err := file.Open()
// 		if err != nil {
// 			return err
// 		}
// 		defer fileReader.Close()

// 		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
// 		if err != nil {
// 			return err
// 		}
// 		defer targetFile.Close()

// 		if _, err := io.Copy(targetFile, fileReader); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func unzipWeb(archive io.ReaderAt, size int64, target string) error {
	reader, err := zip.NewReader(archive, size)
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
