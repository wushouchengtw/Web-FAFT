package main

import (
	"backend/lib"
	"backend/lib/repos"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := lib.GetConfiguration(lib.ConfigurationPath)
	if err != nil {
		log.Fatalf("Can't load the config (%v), get an {%v}", lib.ConfigurationPath, err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", config.Application.Host, config.Application.Port))
	if err != nil {
		log.Fatalf("Can't listen to {%v}:{%v}, got an {%v}", config.Application.Host, config.Application.Port, err)
	}

	db, err := lib.Connect(&config.Database)
	if err != nil {
		log.Fatalf("Can't connect to database")
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("Can't close database connection")
		}
	}(db)

	srv, err := lib.Run(*config, listener)
	if err != nil {
		log.Fatalf("Can't run the server, got an {%v}", err)
	}

	dutRepo := repos.NewDUTRepoInMem()
	testRepo := repos.NewTestRepoInMem()
	resultRepo := repos.NewResultRepoInMem()

	lib.SaveStainlessData("data/20230418-20230424.csv", dutRepo, testRepo, resultRepo)

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
