package main

import (
	"go-final-project/pkg/api"
	"go-final-project/pkg/db"
	"go-final-project/server"
	"log"
	"os"
)

const portDefault int = 7540
const dbDefault string = "scheduler.db"

func main() {
	api.CheckPassword()

	logger := log.New(os.Stdout, "server: ", log.LstdFlags)

	err := db.Init(dbDefault, logger)
	if err != nil {
		logger.Fatal("FATAL: error while db load: ", err)
	}
	defer db.Database.Close()

	srv := server.StartServer(portDefault, logger)
	if err := srv.HTTPServer.ListenAndServe(); err != nil {
		logger.Fatal("FATAL: error while server start: ", err)
	} else {
		logger.Printf("INFO: the server is running")
	}

}
