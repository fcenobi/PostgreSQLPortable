//go:generate goversioninfo -icon=postgresql.ico -manifest=PostgreSQLPortable.manifest
// +build windows
// +build 386 amd64

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const (
	configFile = "PostgreSQLPortable.json"

	strExit            = "Exit"
	strHelp            = "Use the context menu to manage."
	strInit            = "Initializing..."
	strInitFinished    = "Initializing finished"
	strStart           = "Start PostgreSQL Server"
	strStarted         = "Server started"
	strStarting        = "Starting psql shell"
	strStartShell      = "Start PostgreSQL Shell"
	strStartupFinished = "Startup check finished"
	strStop            = "Stop PostgreSQL Server"
	strStopped         = "Server Stopped"
	strTitle           = "PostgreSQL Portable"
)

var (
	ps           = string(filepath.Separator)
	dir, _       = filepath.Abs(filepath.Dir(os.Args[0]))
	pgsqlBaseDir = filepath.Join(dir, "pgsql")
	downloadDir  = filepath.Join(dir, "pgsql-downloads")
	username     = "postgres"

	pgInitdb, pgCtl, pgShell, dataDir, pgHba, logDir, logFile, logPsqlFile     string
	cmdInitDbArgs, cmdStartArgs, cmdStopArgs, cmdStatusArgs, cmdStartShellArgs []string

	serverStatus = false
	serverPid    int

	osName, osArch string
	archiveType    string
	conf           *Configuration
)

func init() {
	checkOs()
	checkArch()
	checkArchiveType()
	conf = new(Configuration)
	loadConfig()
	setPaths()
}

func main() {
	createTray()
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
		showMessage("Please select version first")
		go RunSettingsDialog()
	}
}

func startupCheck() {
	checkNewestVersion()
	findLatest()
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0700)
	}

	if checkStatus() {
		if _, err := os.Stat(pgHba); os.IsNotExist(err) {
			ni.ShowCustom(strTitle, strInit)
			statusPgServerAction.SetText(strInit)
			initdb := exec.Command(pgInitdb, cmdInitDbArgs...)
			initdb.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

			initdbErr := initdb.Run()
			if initdbErr != nil {
				if os.IsNotExist(initdbErr) {
					log.Printf("PostgreSQL files for version %s does not exists!\n", conf.UsedVersion)
				} else {
					log.Printf("initdb error - %s\n", err.Error())
				}
				return
			}
		}
		ni.ShowCustom(strTitle, strInitFinished)
		statusPgServerAction.SetText(strStopped)

		checkExistingVersions()
		log.Println(strStartupFinished)
	}
}

func startingPg() {
	startPg := exec.Command(pgCtl, cmdStartArgs...)
	startPg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	startPgErr := startPg.Run()
	if startPgErr != nil {
		log.Printf("Starting error: %s\n", startPgErr.Error())
	} else {
		serverPid = startPg.Process.Pid
		if serverPid > 0 {
			serverStatus = true
			startPgServerAction.SetEnabled(!serverStatus)
			stopPgServerAction.SetEnabled(serverStatus)
			log.Println(strStarted)
			statusPgServerAction.SetText(strStarted)
			startPgShellAction.SetEnabled(serverStatus)
		}
		checkStatus()
	}
}

func stoppingPg() {
	if serverStatus {
		stopPg := exec.Command(pgCtl, cmdStopArgs...)
		log.Printf("stopPg args - %v\n", stopPg.Args)
		stopPg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		stopPgErr := stopPg.Run()
		if stopPgErr != nil {
			log.Printf("Stopping error: %s\n", stopPgErr.Error())
		} else {
			serverPid = 0
			serverStatus = false
			startPgServerAction.SetEnabled(!serverStatus)
			stopPgServerAction.SetEnabled(serverStatus)
			log.Println(strStopped)
			statusPgServerAction.SetText(strStopped)
			startPgShellAction.SetEnabled(serverStatus)
			checkStatus()
		}
	}
}

func startShell() {
	if serverStatus {
		startSh := exec.Command("cmd", cmdStartShellArgs...)
		startSh.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		startShErr := startSh.Run()
		log.Println(strStarting)
		checkErr("Starting shell error", startShErr)
	}
}

func checkStatus() bool {
	result := false
	if checkExecExists(pgCtl) {
		statusPg := exec.Command(pgCtl, cmdStatusArgs...)
		statusPg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		statusPgErr := statusPg.Run()
		if statusPgErr != nil {
			if statusPgErr.Error() == "exit status 3" {
				log.Println("Server not running!")
			} else {
				log.Printf("checkStatus error: '%s'", statusPgErr.Error())
			}
			result = false
		}
		result = true
	} else {
		log.Printf("Version %s is not installed! Install it first!\n", conf.UsedVersion)
		result = false
	}
	return result
}
