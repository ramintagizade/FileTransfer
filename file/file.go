package file

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var ZipFileName = "files.zip"

func GetFileInfo(FilePath string) (string, error) {
	FileInfo, err := os.Stat(FilePath)

	if err != nil {
		return "", err
	}

	switch mode := FileInfo.Mode(); {
	case mode.IsRegular():
		return "file", nil
	case mode.IsDir():
		return "directory", nil
	}

	return "", nil
}

func ProcessFile(w http.ResponseWriter, r *http.Request, FilePath string) {
	fmt.Println("Processing file .. ", FilePath)
	File, err := os.Open(FilePath)
	if err != nil {
		http.Error(w, "File not found . ", 404)
	}
	defer File.Close()
	FileHeader := make([]byte, 512)
	File.Read(FileHeader)
	FileType := http.DetectContentType(FileHeader)
	FileStat, _ := File.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)
	w.Header().Set("Content-Disposition", "attachment; filename="+FilePath)
	w.Header().Set("Content-Type", FileType)
	w.Header().Set("Content-Length", FileSize)
	File.Seek(0, 0)
	io.Copy(w, File)
	return
}

func ProcessDirectory(w http.ResponseWriter, r *http.Request, FilePath string) {
	files, err := ioutil.ReadDir(FilePath)
	if err != nil {
		return
	}
	zfiles := make([]string, 0)
	for _, file := range files {
		filePath := FilePath + file.Name()
		fi, err := GetFileInfo(filePath)
		if err != nil {
			continue
		}
		switch fi {
		case "file":
			zfiles = append(zfiles, filePath)
		}

	}
	output := ZipFileName
	err = ZipFiles(output, zfiles)
	if err != nil {
		fmt.Println(err)
	}
	http.ServeFile(w, r, output)
}

func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// Get the file information
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Using FileInfoHeader() above only uses the basename of the file. If we want
		// to preserve the folder structure we can overwrite this with the full path.
		header.Name = file

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}

func RemoveZippedFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		return
	}
}
