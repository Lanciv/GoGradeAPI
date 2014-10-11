package main

import (
	"flag"
	"os"
	"strings"

	h "github.com/Lanciv/GoGradeAPI/handlers"
	"github.com/Lanciv/GoGradeAPI/store"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

var log = logrus.New()

var (
	listenAddr     string
	address        string
	dbName         string
	staticDir      string
	insertTestData bool
)

func init() {
	log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = new(logrus.TextFormatter) // default
}

func main() {
	flag.StringVar(&listenAddr, "listenAddr", ":5005", "")
	flag.StringVar(&dbName, "dbName", "dev_go_grade", "")
	flag.BoolVar(&insertTestData, "insertTestData", false, "")

	address = os.Getenv("RETHINKDB_PORT_28015_TCP")
	address = strings.Trim(address, "tcp://")

	if address == "" {
		address = "localhost:28015"
	}

	flag.Parse()

	if err := store.Connect(address, dbName); err != nil {
		log.Fatal("Error setting up database: ", err)
	}

	store.SetupDB(insertTestData)

	r := gin.Default()

	h.SetupHandlers(r)

	r.Run(listenAddr)

}
