// +build windows

package main

import (
	"log"
	"os"

	"github.com/lxn/walk"
)

var (
	mw                   *walk.MainWindow
	ni                   *walk.NotifyIcon
	statusPgServerAction *walk.Action
	startPgServerAction  *walk.Action
	stopPgServerAction   *walk.Action
	startPgShellAction   *walk.Action
	showSettingsAction   *walk.Action

	err error
)

func createTray() {
	mw, err = walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}
	mw.Hide()

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

	// Start status item {{{
	statusPgServerAction = walk.NewAction()
	if err := statusPgServerAction.SetText(strStopped); err != nil {
		log.Fatal(err)
	}
	statusPgServerAction.SetEnabled(false)
	if err := ni.ContextMenu().Actions().Add(statusPgServerAction); err != nil {
		log.Fatal(err)
	}
	// }}} End status item

	// Start separator item {{{
	if err := ni.ContextMenu().Actions().Add(walk.NewSeparatorAction()); err != nil {
		log.Fatal(err)
	}
	// }}} Start separator item

	// Start start item {{{
	startPgServerAction = walk.NewAction()
	if err := startPgServerAction.SetText(strStart); err != nil {
		log.Fatal(err)
	}
	startPgServerAction.SetEnabled(!serverStatus && len(conf.UsedVersion) > 0)

	startPgServerAction.Triggered().Attach(func() { go startingPg() })
	if err := ni.ContextMenu().Actions().Add(startPgServerAction); err != nil {
		log.Fatal(err)
	}
	// }}} End start item

	// Start stop item {{{
	stopPgServerAction = walk.NewAction()
	if err := stopPgServerAction.SetText(strStop); err != nil {
		log.Fatal(err)
	}
	stopPgServerAction.SetEnabled(serverStatus)

	stopPgServerAction.Triggered().Attach(func() { go stoppingPg() })
	if err := ni.ContextMenu().Actions().Add(stopPgServerAction); err != nil {
		log.Fatal(err)
	}
	// }}} End stop item

	if err := ni.ContextMenu().Actions().Add(walk.NewSeparatorAction()); err != nil {
		log.Fatal(err)
	}

	// Start shell item {{{
	startPgShellAction = walk.NewAction()
	if err := startPgShellAction.SetText(strStartShell); err != nil {
		log.Fatal(err)
	}
	startPgShellAction.SetEnabled(false)
	startPgShellAction.Triggered().Attach(func() { go startShell() })
	if err := ni.ContextMenu().Actions().Add(startPgShellAction); err != nil {
		log.Fatal(err)
	}
	// }}} End shell item

	if err := ni.ContextMenu().Actions().Add(walk.NewSeparatorAction()); err != nil {
		log.Fatal(err)
	}

	// Start settings item {{{
	showSettingsAction = walk.NewAction()
	if err := showSettingsAction.SetText("Settings"); err != nil {
		log.Fatal(err)
	}
	showSettingsAction.Triggered().Attach(func() { RunSettingsDialog() })
	if err := ni.ContextMenu().Actions().Add(showSettingsAction); err != nil {
		log.Fatal(err)
	}
	// }}} End settings item

	// Start exit item {{{
	exitAction := walk.NewAction()
	if err := exitAction.SetText(strExit); err != nil {
		log.Fatal(err)
	}

	exitAction.Triggered().Attach(func() { quit() })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}
	// }}} End exit item

	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	mw.Run()
}

func appQuit() {
	// walk.App().Exit(0)
	os.Exit(0)
}
