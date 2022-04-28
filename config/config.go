package config

import (
	"github.com/ShevchenkoVadim/helperlib/sfotypes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var C = sfotypes.Configurations{}

func init() {
	log.Println("Load config file...")
	configName := filepath.Base(os.Args[0])
	f, err := ioutil.ReadFile(configName + ".yaml")

	if err != nil {
		log.Fatal("Cannot open config file")
		//os.Exit(101)
		return
	}

	err = yaml.Unmarshal(f, &C)
	if err != nil {
		log.Fatal("Cannot read config file")
		//os.Exit(102)
		return
	}
	log.Println(C)
}
