//go:generate goversioninfo -icon=postgresql.ico -manifest=PostgreSQLPortable.manifest
// +build windows

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/lxn/walk"
)

var (
	dir, _      = filepath.Abs(filepath.Dir(os.Args[0]))
	pgInitdb    = filepath.Join(dir, "/bin/initdb")
	pgCtl       = filepath.Join(dir, "/bin/pg_ctl")
	pgShell     = filepath.Join(dir, "/bin/psql")
	dataDir     = filepath.Join(dir, "/data")
	pgHba       = filepath.Join(dir, "/data/pg_hba.conf")
	logDir      = filepath.Join(dir, "/log")
	logFile     = filepath.Join(logDir, "/postgres.log")
	logPsqlFile = filepath.Join(logDir, "/psql.log")
	username    = "postgres"

	strExit            = "E&xit"
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

	cmdInitDbArgs     = []string{"-D", dataDir, "-U", username, "-A", "trust", "-E", "UTF8", "--locale=russian_russia", "-k", "-n"}
	cmdStartArgs      = []string{"-D", dataDir, "-l", logFile, "-w", "start"}
	cmdStopArgs       = []string{"-D", dataDir, "stop"}
	cmdStatusArgs     = []string{"-D", dataDir, "status"}
	cmdStartShellArgs = []string{"/C", "start", "/wait", pgShell, "-L", logPsqlFile, "-U", username, username}

	serverStatus         = false
	serverPid            int
	ni                   *walk.NotifyIcon
	statusPgServerAction *walk.Action
	startPgServerAction  *walk.Action
	stopPgServerAction   *walk.Action
	startPgShellAction   *walk.Action
)

func main() {
	mw, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}

	icon, err := walk.NewIconFromResourceId(10)
	if err != nil {
		log.Fatal(err)
	}

	ni, err = walk.NewNotifyIcon()
	if err != nil {
		log.Fatal(err)
	}
	defer ni.Dispose()

	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetToolTip(strTitle); err != nil {
		log.Fatal(err)
	}

	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}

		if err := ni.ShowCustom(strTitle, strHelp); err != nil {
			log.Fatal(err)
		}
	})

	statusPgServerAction = walk.NewAction()
	if err := statusPgServerAction.SetText(strStopped); err != nil {
		log.Fatal(err)
	}
	statusPgServerAction.SetEnabled(false)
	if err := ni.ContextMenu().Actions().Add(statusPgServerAction); err != nil {
		log.Fatal(err)
	}

	separatorAction := walk.NewAction()
	if err := separatorAction.SetText("-"); err != nil {
		log.Fatal(err)
	}
	separatorAction.SetEnabled(false)
	if err := ni.ContextMenu().Actions().Add(separatorAction); err != nil {
		log.Fatal(err)
	}

	startPgServerAction = walk.NewAction()
	if err := startPgServerAction.SetText(strStart); err != nil {
		log.Fatal(err)
	}
	startPgServerAction.SetEnabled(!serverStatus)

	startPgServerAction.Triggered().Attach(func() { go startingPg() })
	if err := ni.ContextMenu().Actions().Add(startPgServerAction); err != nil {
		log.Fatal(err)
	}

	stopPgServerAction = walk.NewAction()
	if err := stopPgServerAction.SetText(strStop); err != nil {
		log.Fatal(err)
	}
	stopPgServerAction.SetEnabled(serverStatus)

	stopPgServerAction.Triggered().Attach(func() { go stopingPg() })
	if err := ni.ContextMenu().Actions().Add(stopPgServerAction); err != nil {
		log.Fatal(err)
	}

	separatorAction2 := walk.NewAction()
	if err := separatorAction2.SetText("-"); err != nil {
		log.Fatal(err)
	}
	separatorAction2.SetEnabled(false)
	if err := ni.ContextMenu().Actions().Add(separatorAction2); err != nil {
		log.Fatal(err)
	}

	startPgShellAction = walk.NewAction()
	if err := startPgShellAction.SetText(strStartShell); err != nil {
		log.Fatal(err)
	}
	startPgShellAction.SetEnabled(false)

	startPgShellAction.Triggered().Attach(func() { go startShell() })
	if err := ni.ContextMenu().Actions().Add(startPgShellAction); err != nil {
		log.Fatal(err)
	}

	separatorAction3 := walk.NewAction()
	if err := separatorAction3.SetText("-"); err != nil {
		log.Fatal(err)
	}
	separatorAction3.SetEnabled(false)
	if err := ni.ContextMenu().Actions().Add(separatorAction3); err != nil {
		log.Fatal(err)
	}

	exitAction := walk.NewAction()
	if err := exitAction.SetText(strExit); err != nil {
		log.Fatal(err)
	}

	exitAction.Triggered().Attach(func() { stopingPg(); walk.App().Exit(0) })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}

	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	go startupCheck()

	mw.Run()
}

func startupCheck() {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0700)
	}

	if _, err := os.Stat(pgHba); os.IsNotExist(err) {
		ni.ShowCustom(strTitle, strInit)
		statusPgServerAction.SetText(strInit)
		initdb := exec.Command(pgInitdb, cmdInitDbArgs...)
		initdb.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		initdbErr := initdb.Run()
		if initdbErr != nil {
			log.Println(err.Error())
			return
		}
	}
	ni.ShowCustom(strTitle, strInitFinished)
	statusPgServerAction.SetText(strStopped)

	log.Println(strStartupFinished)
}

func startingPg() {
	startPg := exec.Command(pgCtl, cmdStartArgs...)
	startPg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	startPgErr := startPg.Run()
	if startPgErr != nil {
		log.Printf("Starting error: %s\n", startPgErr.Error())
	}
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

func stopingPg() {
	stopPg := exec.Command(pgCtl, cmdStopArgs...)
	stopPg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	stopPgErr := stopPg.Run()
	if stopPgErr != nil {
		log.Printf("Stopping error: %s\n", stopPgErr.Error())
	}
	serverPid = 0
	serverStatus = false
	startPgServerAction.SetEnabled(!serverStatus)
	stopPgServerAction.SetEnabled(serverStatus)
	log.Println(strStopped)
	statusPgServerAction.SetText(strStopped)
	startPgShellAction.SetEnabled(serverStatus)
	checkStatus()
}

func startShell() {
	if serverStatus {
		startSh := exec.Command("cmd", cmdStartShellArgs...)
		startSh.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		startShErr := startSh.Run()
		log.Println(strStarting)
		if startShErr != nil {
			log.Printf("Starting shell error: %s\n", startShErr.Error())
		}
	}
}

func checkStatus() {
	statusPg := exec.Command(pgCtl, cmdStatusArgs...)
	statusPg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, statusPgErr := statusPg.Output()
	if statusPgErr != nil {
		log.Printf("checkStatus error: %s", statusPgErr)
	}
	log.Println(string(out))
}
