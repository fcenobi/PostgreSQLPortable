package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

func setPaths() {
	if len(conf.UsedVersion) > 0 {
		pgInitdb = filepath.Join(pgsqlBaseDir, conf.UsedVersion, "bin/initdb")
		pgCtl = filepath.Join(pgsqlBaseDir, conf.UsedVersion, "bin/pg_ctl")
		pgShell = filepath.Join(pgsqlBaseDir, conf.UsedVersion, "bin/psql")
		dataDir = filepath.Join(pgsqlBaseDir, conf.UsedVersion, "data")
		pgHba = filepath.Join(pgsqlBaseDir, conf.UsedVersion, "data/pg_hba.conf")
		logDir = filepath.Join(pgsqlBaseDir, conf.UsedVersion, "log")
		logFile = filepath.Join(logDir, "postgres.log")
		logPsqlFile = filepath.Join(logDir, "psql.log")

		cmdInitDbArgs = []string{"-D", dataDir, "-U", username, "-A", "trust", "-E", "UTF8", "--locale=american_usa", "-k", "-n"}
		cmdStartArgs = []string{"-D", dataDir, "-l", logFile, "-w", "start"}
		cmdStopArgs = []string{"-D", dataDir, "stop"}
		cmdStatusArgs = []string{"-D", dataDir, "status"}
		cmdStartShellArgs = []string{"/C", "start", "/wait", pgShell, "-L", logPsqlFile, "-U", username, username}

	} else {
		ShowMessage(strPSVF)
		go ShowSettingsDialog()
	}
}

//func dataExists() {
//	if checkServerStatus() {
//		if _, err := os.Stat(pgHba); os.IsNotExist(err) {
//			ShowNotification(strInit)
//			SetStatus(strInit)
//			initdb()
//		}
//		ShowNotification(strInitFinished)
//		SetStatus(strStopped)
//		log.Println(strStartupFinished)
//	}
//}

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
		ShowNotification(fmt.Sprintf(strDVPW, ver))
		SetNotificationToolTip(fmt.Sprintf(strDVPW, ver))

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
		ShowNotification("Download finished!")
	} else {
		log.Println("Already downloaded")
	}
	return destName
}

func install(v string) {
	SetStatus(fmt.Sprintf(strInstalling, v))
	EnableStart(false)

	f := download(v)
	extract(f, v)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0700)
	}
	initdb()
	ShowNotification("")
	checkServerStatus()
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
	stopPg()
	AppQuit()
}
