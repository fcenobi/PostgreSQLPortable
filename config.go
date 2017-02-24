package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Configuration struct {
	UsedVersion       string
	CheckForUpdates   bool
	AutoInstallLatest bool
	Username          string
	Locale            string
}

func loadConfig() {
	log.Println("Reading config...")
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.SetPrefix("ERROR")
		log.Printf("Can't read config file : %s\n", err.Error())
		return
	}
	err = json.Unmarshal(bytes, &conf)
	checkErr("Can't parse config", err)
}

func saveConfig() error {
	log.Println("Writing config...")
	bytes, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFile, bytes, 0644)
}
