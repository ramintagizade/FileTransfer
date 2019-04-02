package server

import (
	. "../file"
	. "../utils"
	"fmt"
	"net/http"
	"strings"
)

var mapFile = make(map[string]string)

func Handler(w http.ResponseWriter, r *http.Request) {
	FileName := mapFile[strings.Replace(r.URL.Path, "/", "", -1)]
	fi, err := GetFileInfo(FileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch fi {
	case "file":
		ProcessFile(w, r, FileName)
	case "directory":
		ProcessDirectory(w, r, FileName)
	}

	RemoveZippedFile(ZipFileName)
	return
}

func Run(filepath string) {
	pkgLink := RandomString(5)
	mapFile[pkgLink] = filepath
	ipAddress := GetAddress()
	port := ":8000"
	downloadLink := "http://" + ipAddress + ":" + port + "/" + pkgLink
	fmt.Println("download link: ", downloadLink)
	http.HandleFunc("/"+pkgLink, Handler)
	http.ListenAndServe(port, nil)
}
