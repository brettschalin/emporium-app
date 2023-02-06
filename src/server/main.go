package main

import (
	"context"
	"fmt"
	"os"

	"marketplace-server/config"
	"marketplace-server/db"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New("./config/config.json")
	if err != nil {
		fmt.Printf("Could not load config: %s\n", err)
		os.Exit(1)
	}

	err = db.Connect(ctx, cfg.DatabaseConnString())
	if err != nil {
		fmt.Printf("Could not connect to DB: %s\n", err)
		os.Exit(1)
	}

	defer db.Cleanup()

	r := gin.Default()
	r.SetTrustedProxies(nil)

	SetupRoutes(r)

	r.Run()
}
