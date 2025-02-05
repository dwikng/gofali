package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dwikng/gofali/storage"
	"github.com/vharitonsky/iniflags"
)

type config struct {
	host string
	port int

	mysqlConnect string

	adminPath    string
	rootRedirect string
	slugLength   int
}

var conf config

func main() {
	flag.StringVar(&conf.host, "host", "127.0.0.1", "server host")
	flag.IntVar(&conf.port, "port", 8080, "server port")
	flag.StringVar(&conf.mysqlConnect, "mysql-connect", "gofali:gofali@tcp(127.0.0.1:3306)/gofali", "mysql connection string")
	flag.StringVar(&conf.adminPath, "admin-path", "admin", "admin path")
	flag.StringVar(&conf.rootRedirect, "root-redirect", "", "redirect path")
	flag.IntVar(&conf.slugLength, "slug-length", 8, "slug length")
	iniflags.Parse()

	store, err := storage.NewMysql(conf.mysqlConnect, conf.slugLength)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	a := &api{
		store: store,
	}

	go func() {
		if err = a.run(); err != nil {
			log.Fatalf("error starting fasthttp server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		if err = a.shutdown(); err != nil {
			log.Fatalf("error shutting down fasthttp server: %v", err)
		}
		if err = store.Close(); err != nil {
			log.Fatalf("error closing database: %v", err)
		}
	}
}
