package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Configuration struct {
	UsedVersion       string
	CheckForUpdates   bool
	AutoInstallLatest bool
	Username          string
	Locale            string
}

func NewConfiguration() *Configuration {
	return &Configuration{UsedVersion: "", CheckForUpdates: false, AutoInstallLatest: false, Username: "postgres", Locale: "american_usa"}
}

func loadConfig() error {
	bytes, err := ioutil.ReadFile(filepath.Join(dir, configFile))
	if err != nil {
		log.Printf("Can't read config file: %s\n", err.Error())
		return err
	}
	if err = json.Unmarshal(bytes, &conf); err != nil {
    	log.Printf("Can't parse config: %s\n", err.Error())
    	return err
    }
	return nil
}

func saveConfig() error {
	bytes, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(dir, configFile), bytes, 0644)
}
