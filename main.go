package main

import (
	"bufio"
	"flag"
	"github.com/gr4y/fritzbox-call-monitor/data"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net"
	"os"
	"os/signal"
)

var (
	router       = flag.String("router", "fritz.box:1012", "Hostname and Port to TCP Port on fritz.box")
	logFile      = flag.String("logFile", "/var/log/fritzbox-call-monitor.log", "Path to log file")
	databaseFile = flag.String("databaseFile", "data.db", "Path to database file")
	database     gorm.DB
)

const (
	STRING_DELIMITER = ";"
	LINE_DELIMITER   = '\n'
)

func init() {
	flag.Parse()
}

func main() {
	database := openDatabaseConnection()
	database.LogMode(true)
	if !database.HasTable(&data.CallEvent{}) {
		database.CreateTable(&data.CallEvent{})
	}

	conn, err := net.Dial("tcp", *router)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	buffer := bufio.NewReader(conn)
	for {
		str, err := buffer.ReadString(LINE_DELIMITER)
		if err != nil {
			log.Println(err)
		}
		event := data.NewEvent(str)
		if event != nil && event.IsValid() {
			if database.NewRecord(event) {
				database.Create(event)
			}
		}
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			log.Println("Received an Interrupt. Shutting down monitor")
			conn.Close()
			database.Close()
		}
	}()
	<-cleanupDone

}

func openDatabaseConnection() gorm.DB {
	db, err := gorm.Open("sqlite3", *databaseFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return db
}
