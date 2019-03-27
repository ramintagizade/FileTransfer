package server

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var mapFile = make(map[string]string)

func Handler(w http.ResponseWriter, r *http.Request) {
	FileName := mapFile[strings.Replace(r.URL.Path, "/", "", -1)]
	File, err := os.Open(FileName)
	if err != nil {
		http.Error(w, "File not found . ", 404)
		return
	}
	defer File.Close()

	FileHeader := make([]byte, 512)
	File.Read(FileHeader)
	FileType := http.DetectContentType(FileHeader)
	FileStat, _ := File.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)
	w.Header().Set("Content-Disposition", "attachment; filename="+FileName)
	w.Header().Set("Content-Type", FileType)
	w.Header().Set("Content-Length", FileSize)
	File.Seek(0, 0)
	io.Copy(w, File)
	return
}

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	s := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = s[rand.Intn(len(s))]
	}
	return string(b)
}

func Run(filepath string) {
	pkgLink := randomString(5)
	mapFile[pkgLink] = filepath
	ipAddress := getAddress()
	port := "8000"
	downloadLink := "http://" + ipAddress + ":" + port + "/" + pkgLink
	fmt.Println("download link: ", downloadLink)
	http.HandleFunc("/"+pkgLink, Handler)
	http.ListenAndServe(":"+port, nil)
	fmt.Println("Server is running on port " + port)
}

func getAddress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Get address : ", err)
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println("Error getting address")
			continue
		}

		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPAddr:
				ipad := v.IP.To4().String()
				if strings.Contains(ipad, "192") {
					return ipad
				}
			case *net.IPNet:
				ipad := v.IP.To4().String()
				if strings.Contains(ipad, "192") {
					return ipad
				}
			}

		}
	}
	return ""
}
