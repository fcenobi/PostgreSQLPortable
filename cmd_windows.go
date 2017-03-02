package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func initdb() bool {
	initdb := exec.Command(pgInitdb, cmdInitDbArgs...)
	initdb.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	initdbErr := initdb.Run()
	if initdbErr != nil {
		if os.IsNotExist(initdbErr) {
			log.Printf(strFNEErr, conf.UsedVersion)
		} else {
			log.Printf("initdb error - %s\n", err.Error())
		}
		return false
	}
	return true
}

func startPg() {
	if checkServerStatus() {
		startPg := exec.Command(pgCtl, cmdStartArgs...)
		startPg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		startPgErr := startPg.Run()
		if startPgErr != nil {
			log.Printf(strStartErr, startPgErr.Error())
			log.Println("Server not started!")
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
		}
	} else {
		log.Println("Server not started!")
	}
}

func stopPg() {
	if serverStatus {
		stopPg := exec.Command(pgCtl, cmdStopArgs...)
		stopPg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		stopPgErr := stopPg.Run()
		if stopPgErr != nil {
			log.Printf(strStopErr, stopPgErr.Error())
		} else {
			serverPid = 0
			serverStatus = false
			startPgServerAction.SetEnabled(!serverStatus)
			stopPgServerAction.SetEnabled(serverStatus)
			log.Println(strStopped)
			statusPgServerAction.SetText(strStopped)
			startPgShellAction.SetEnabled(serverStatus)
			checkServerStatus()
		}
	}
}

func checkServerStatus() bool {
	result := false
	if checkExecExists(pgCtl) {
		statusPg := exec.Command(pgCtl, cmdStatusArgs...)
		statusPg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		statusPgErr := statusPg.Run()
		if statusPgErr != nil {
			if statusPgErr.Error() == "exit status 3" {
				log.Println(strSNR)
			} else if statusPgErr.Error() == "exit status 4" {
				log.Println(strDNI)
				answer := AskQuestion(fmt.Sprintf("%s Initialize?", strDNI))
				if answer == 1 {
					return initdb()
				} else {
					return false
				}
			} else {
				log.Printf("checkStatus error: '%s'", statusPgErr.Error())
			}
			result = false
		}
		result = true
	} else {
		log.Printf(strIIF, conf.UsedVersion)
		result = false
	}
	return result
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
