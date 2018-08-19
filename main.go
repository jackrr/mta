package main

import (
	"fmt"
	"github.com/jackrr/mta/api"
	"gopkg.in/ini.v1"
	"os"
)

func main() {
	cfg, err := ini.Load("api.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	m := api.MTA{ApiKey: cfg.Section("").Key("api_key").String()}
	m.GetData()
}
