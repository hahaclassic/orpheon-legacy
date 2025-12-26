package main

import (
	"flag"

	"github.com/hahaclassic/orpheon/backend/internal/app"
	"github.com/hahaclassic/orpheon/backend/internal/config"
)

var configPath = ".env"

func init() {
	flag.StringVar(&configPath, "config", ".env", "path to config file")
	flag.Parse()
}

func main() {
	conf := config.MustLoad(configPath)

	app.Run(conf)
}
