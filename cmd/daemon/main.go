package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/daos/daos/pkg/api"
	"github.com/daos/daos/pkg/config"
	"github.com/daos/daos/pkg/db"
	"github.com/daos/daos/cmd/daemon/handlers"
	"github.com/gin-gonic/gin"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	configPath := flag.String("config", "/opt/daos/config.yaml", "Path to config file")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("DAOS Daemon %s (commit: %s, date: %s)\n", version, commit, date)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Printf("Failed to load config from %s, using defaults: %v", *configPath, err)
		cfg = config.Default()
	}

	database, err := db.New(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	router := gin.Default()
	api.RegisterHandlers(router, handlers.New(database))

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	
	go func() {
		log.Printf("Starting DAOS daemon on %s", addr)
		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down DAOS daemon...")
}
