package lib

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func engine(config *Configuration) *gin.Engine {
	engine := gin.Default()

	mode := strings.ToLower(config.Application.Mode)
	switch mode {
	case "test":
		gin.SetMode(gin.TestMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.TestMode)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	engine.Use(cors.New(corsConfig))

	return engine
}

func Run(config Configuration, listener net.Listener) (*http.Server, error) {
	engine := engine(&config)

	srv := &http.Server{Handler: engine.Handler()}

	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Can't create a server, get an {%v}", err)
		}
	}()

	return srv, nil
}
