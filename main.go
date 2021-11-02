package main

import (
	_ "database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gcfg.v1"
	"log"
	"oauth2-gmail/database"
	"oauth2-gmail/model"
	"oauth2-gmail/server"
)

func main() {
	success := database.InitTable()
	if !success {
		log.Fatal("CreateUsersTbl false")
	}
	model.GlbConfig = model.Config{}
	err := gcfg.ReadFileInto(&model.GlbConfig, "./template.conf")

	if err != nil {
		log.Fatal(err.Error())
	}
	go server.StartExtServer(model.GlbConfig)
	server.StartIntServer(model.GlbConfig)
	fmt.Println(model.GlbConfig)
}
