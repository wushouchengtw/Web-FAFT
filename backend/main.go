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
)

func main() {
	config, err := lib.GetConfiguration(lib.ConfigurationPath)
	if err != nil {
		log.Fatalf("Can't load the config (%v), get an {%v}", lib.ConfigurationPath, err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", config.Host, config.Port))
	if err != nil {
		log.Fatalf("Can't listen to {%v}:{%v}, got an {%v}", config.Host, config.Port, err)
	}

	srv, err := lib.Run(*config, listener)
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
