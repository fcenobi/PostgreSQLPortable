package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func getMax(x []string) string {
	var max string
	if len(x) == 0 {
		log.Println("Array must be non empty")
		return ""
	} else {
		max = x[0]
		for _, v := range x {
			if v > max {
				max = v
			}
		}
	}
	return max
}

func checkOs() {
	switch runtime.GOOS {
	case "windows":
		osName = "windows"
	case "linux":
		osName = "linux"
	case "darwin":
		osName = "osx"
	default:
		log.Printf("Unknown Platform %s\n", runtime.GOOS)
		os.Exit(1)
	}
}

func checkArch() {
	switch runtime.GOARCH {
	case "386":
		osArch = ""
	case "amd64":
		osArch = "-x64"
	default:
		log.Printf("Unknown architecture %s\n", runtime.GOARCH)
		os.Exit(1)
	}
}

func checkArchiveType() {
	switch osName {
	case "linux":
		archiveType = "tar.gz"
	case "windows", "osx":
		archiveType = "zip"
	default:
		log.Println("Unknown archive type")
		os.Exit(1)
	}
}

func checkErr(msg string, err error) {
	if err != nil {
		log.Printf("ERROR: %s - %s\n", msg, err.Error())
		os.Exit(1)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func download(ver string) string {
	url := fmt.Sprintf("http://get.enterprisedb.com/postgresql/postgresql-%s-1-%s%s-binaries.%s", ver, osName, osArch, archiveType)
	tokens := strings.Split(url, "/")
	destDir := downloadDir
	fileName := tokens[len(tokens)-1]

	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		os.MkdirAll(destDir, 0755)
	}
	destName := fmt.Sprintf("%s/%s", destDir, fileName)

	if _, err := os.Stat(destName); os.IsNotExist(err) {
		ni.ShowCustom(strTitle, fmt.Sprintf("Downloading version %s", ver))
		ni.SetToolTip(fmt.Sprintf("%s: Downloading version %s. Please wait", strTitle, ver))

		log.Printf("Downloading %s to %s\n", url, destName)

		output, err := os.Create(destName)
		if err != nil {
			log.Printf("Error while creating %s - %s\n", fileName, err.Error())
			return ""
		}
		defer output.Close()

		response, err := http.Get(url)
		if err != nil {
			log.Printf("Error while downloading %s - %s\n", url, err.Error())
			return ""
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			log.Printf("Can't download! Status code: %d\n", response.StatusCode)
		}

		dest, err := os.Create(destName)
		if err != nil {
			fmt.Printf("Can't create %s: %v\n", destName, err)
			return ""
		}
		defer dest.Close()

		writer := io.MultiWriter(dest)

		n, err := io.Copy(writer, response.Body)
		if err != nil {
			log.Printf("Error while downloading %s - %s\n", url, err.Error())
			return ""
		}

		log.Println(n, "bytes downloaded.")
		ni.ShowInfo(strTitle, "Download finished!")
	} else {
		log.Printf("Already downloaded")
	}
	return destName
}

func install(v string) {
	f := download(v)
	extract(f, v)
}

func checkExecExists(c string) bool {
	_, err := exec.LookPath(c)
	if err != nil {
		log.Printf("%s is not exist", c)
		return false
	}
	return true
}

func quit() {
	stoppingPg()
	appQuit()
}
