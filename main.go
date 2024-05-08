package main

import (
	"employee_manage_svcs/router"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func initLogger(logDir *string) {
	fPath := ""
	if *logDir != "" {
		fPath = fmt.Sprintf("%vaccess.log", *logDir)
	} else {
		wd, _ := os.Getwd()
		fPath = wd + "-access.log"
	}
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.ToSlash(fPath),
		MaxSize:    1, // MB
		MaxBackups: 10,
		MaxAge:     10, // days
	}

	// Fork writing into two outputs
	multiWriter := io.MultiWriter(os.Stderr, lumberjackLogger)

	logFormatter := new(log.TextFormatter)
	logFormatter.TimestampFormat = time.RFC1123Z // or RFC3339
	logFormatter.FullTimestamp = true

	log.SetFormatter(logFormatter)
	log.SetLevel(log.InfoLevel)
	log.SetOutput(multiWriter)
}

func main() {
	logDir := flag.String("logdir", "", "logging directory")
	flag.Parse()
	initLogger(logDir)
	router := router.InitRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Println("Launching the app, visit localhost:8000/")
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
