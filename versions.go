package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/blang/semver"
)

func checkAvailableVersions() []string {
	var availableVersions []string
	if doc, err := goquery.NewDocument("https://www.postgresql.org/"); err == nil {
		doc.Find("#pgFrontLatestReleasesWrap > b").Each(func(i int, s *goquery.Selection) {
			ver := strings.TrimSpace(s.Text())
			availableVersions = append(availableVersions, ver)
		})
	} else {
		ShowNotification("Can't check available versions")
		log.Println("Can't check available versions")
	}
	return availableVersions
}

func checkExistingVersions() []string {
	var existingVersions []string
	if _, err := os.Stat(pgsqlBaseDir); os.IsNotExist(err) {
		log.Println("Base dir is not exist! Creating...")
	}
	files, _ := ioutil.ReadDir(pgsqlBaseDir)
	for _, f := range files {
		if f.IsDir() {
			existingVersions = append(existingVersions, f.Name())
		}
	}
	if len(existingVersions) == 0 && conf.AutoInstallLatest {
		latest := getMax(checkAvailableVersions())
		go install(latest)
	}
	return existingVersions
}

func checkNewestVersion() {
	var resultMsg string
	ev := checkExistingVersions()
	av := checkAvailableVersions()
	if len(ev) == 0 {
		resultMsg = strNIV
	}
	for _, e := range ev {
		ve, _ := semver.Make(e)
		for _, a := range av {
			if contains(ev, a) {
				continue
			}
			va, _ := semver.Make(a)
			if (ve.Major == va.Major) && (ve.Minor == va.Minor) {
				if ve.Patch < va.Patch {
					resultMsg += fmt.Sprintf("New version '%s' available for %d.%d (existing %s)\n", va, ve.Major, ve.Minor, ve)
				}
			}
		}
	}
	if len(resultMsg) != 0 {
		ShowNotification(resultMsg)
	}
}

func findLatest() {
	log.Printf("Latest existing version - %s\n", getMax(checkExistingVersions()))
	log.Printf("Latest available version - %s\n", getMax(checkAvailableVersions()))
}
