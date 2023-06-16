package main

import (
	"backend/lib"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	config *lib.Configuration
	err    error
)

func init() {
	config, err = lib.GetConfiguration(lib.ConfigurationPath)
	if err != nil {
		log.Fatalf("Can't load the config (%v), get an {%v}", lib.ConfigurationPath, err)
	}

}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", config.Host, config.Port))
	if err != nil {
		log.Fatalf("Can't listen to {%v}:{%v}, got an {%v}", config.Host, config.Port, err)
	}
	// To-do pass a wrong username & password temporary.
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", "f1", "f2", "f3", "f4"))
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}

	srv, err := lib.Run(*config, listener, db)
	if err != nil {
		log.Fatalf("Can't run the server, got an {%v}", err)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	log.Println("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown server got an {%v}", err)
	}

	select {
	case <-ctx.Done():
		log.Println("Shutdown server timeout")
	}

	log.Println("Gracefully shutdown server")
}
