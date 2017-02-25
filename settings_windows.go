package main

import (
	"fmt"
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	settingsDlg = &SettingsDialogWindow{model: NewAVModel()}

	db *walk.DataBinder
	ep walk.ErrorPresenter
)

type SettingsDialogWindow struct {
	*walk.Dialog
	model               *AVModel
	autoInstallLatestCB *walk.CheckBox
	availableVersionsLB *walk.ListBox
	checkForUpdateCB    *walk.CheckBox
	existingVersionsCMB *walk.ComboBox
	usernameLE          *walk.LineEdit
	localeLE            *walk.LineEdit
	acceptPB            *walk.PushButton
	cancelPB            *walk.PushButton
}

type AVItem struct {
	name  string
	value string
}

type AVModel struct {
	walk.ListModelBase
	items []AVItem
}

func NewAVModel() *AVModel {
	av := checkAvailableVersions()

	m := &AVModel{items: make([]AVItem, len(av))}

	for i, v := range av {
		name := v
		value := v
		m.items[i] = AVItem{name, value}
	}
	return m
}

func (m *AVModel) ItemCount() int {
	return len(m.items)
}

func (m *AVModel) Value(index int) interface{} {
	return m.items[index].name
}

func (sd *SettingsDialogWindow) ev_ItemActivated() {
	value := settingsDlg.model.items[settingsDlg.availableVersionsLB.CurrentIndex()].value
	answer := walk.MsgBox(nil, "PostgreSQL Portable", fmt.Sprintf("Are you sure to install PostgreSQL %s?", value), walk.MsgBoxOKCancel+walk.MsgBoxIconQuestion)
	if answer == 1 {
		log.Printf("Installing PostgreSQL %s\n", value)
		go install(value)
	}
}

func showMessage(msg string) {
	walk.MsgBox(nil, strTitle, msg, walk.MsgBoxIconInformation)
}

func RunSettingsDialog() (int, error) {
	SD := Dialog{
		AssignTo: &settingsDlg.Dialog,
		Title:    "Settings",
		DataBinder: DataBinder{
			AssignTo:   &db,
			DataSource: conf,
		},
		Size:      Size{250, 200},
		MinSize:   Size{250, 200},
		MaxSize:   Size{250, 200},
		FixedSize: true,
		Layout:    VBox{Margins{5, 5, 5, 5}, 5, false, false},
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{
					Label{
						Text: "Existing versions:",
					},
					ComboBox{
						AssignTo: &settingsDlg.existingVersionsCMB,
						Value:    Bind("UsedVersion"),
						Model:    checkExistingVersions(),
					},

					Label{
						Text: "Available versions:",
					},
					ListBox{
						AssignTo:        &settingsDlg.availableVersionsLB,
						MinSize:         Size{50, 70},
						MaxSize:         Size{50, 70},
						Model:           settingsDlg.model,
						OnItemActivated: settingsDlg.ev_ItemActivated,
					},
				},
			},
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Username:",
					},
					LineEdit{
						AssignTo:    &settingsDlg.usernameLE,
						Text:        Bind("Username"),
						ToolTipText: "postgres for example",
						CueBanner:   "postgres for example",
					},
					Label{
						Text: "Locale:",
					},
					LineEdit{
						AssignTo:    &settingsDlg.localeLE,
						Text:        Bind("Locale"),
						ToolTipText: "american_usa or russian_russia for example",
						CueBanner:   "american_usa or russian_russia for example",
					},
					Label{
						Text: "Autoinstall latest:",
					},
					CheckBox{
						AssignTo: &settingsDlg.autoInstallLatestCB,
						Checked:  Bind("AutoInstallLatest"),
					},
					Label{
						Text: "Check for update:",
					},
					CheckBox{
						AssignTo: &settingsDlg.checkForUpdateCB,
						Checked:  Bind("CheckForUpdates"),
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &settingsDlg.acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Println(err)
								return
							}
							saveConfig()
							setPaths()
							settingsDlg.Hide()
						},
					},
					PushButton{
						AssignTo: &settingsDlg.cancelPB,
						Text:     "Cancel",
						OnClicked: func() {
							settingsDlg.Hide()
						},
					},
				},
			},
		},
	}
	return SD.Run(nil)
}
