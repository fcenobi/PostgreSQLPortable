// +build windows linux
// +build 386 amd64

package main

import (
	"os"
	"path/filepath"
)

const (
	configFile = "PostgreSQLPortable.json"

	strOK           = "OK"
	strCancel       = "Cancel"
	strExit         = "Exit"
	strCFU          = "Check For Updates"
	strHelp         = "Use the context menu to manage."
	strInit         = "Initializing..."
	strInitFinished = "Initializing finished"
	strPSVF         = "Please select version first"
	strExtrV        = "Extracting version %s"
	strExtrVF       = "Extracting version %s finished"
	strInstalling   = "Installing PostgreSQL %s\n"
	strInstallation = "Installation Finished"
	strSettings     = "Settings"
	strStart        = "Start PostgreSQL Server"
	strStarted      = "Server Started"
	strStarting     = "Starting psql shell"
	strStartShell   = "Start PostgreSQL Shell"
	strStop         = "Stop PostgreSQL Server"
	strStopped      = "Server Stopped"
	strTitle        = "PostgreSQL Portable"
	strSNR          = "Server not running!"
	strDNI          = "Database is not initialized!"
	strNIV          = "No installed versions"
	strDVPW         = "Downloading version %s. Please wait"
	strAUSI         = "Are you sure to install PostgreSQL %s?"
	strIIF          = "Version %s is not installed! Install it first!\n"
	strStopErr      = "Stopping error: %s\n"
	strStartErr     = "Starting error: %s\n"
	strFNEErr       = "PostgreSQL files for version %s does not exists!\n"
)

var (
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
}

func main() {
	conf = NewConfiguration()
	loadConfig()
	setPaths()
	CreateTray()
}
