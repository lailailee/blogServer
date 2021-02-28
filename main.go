package main

import (
	"end/core"
	"end/models"
	"end/router"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
)

func main() {
	core.Config = core.NewConfig()
	var configPath = flag.String("c", "./config/config.json", "help message for flagname")
	flag.Parse()
	if err := core.Config.Loads(*configPath); err != nil {
		fmt.Printf("Load ConfigFile %v, error, %v", configPath, err)
		return
	}

	os.MkdirAll(core.LogPath, os.ModePerm)

	core.InitLogger()

	models.InitDB()

	// start server
	r := gin.Default()
	router.GinRouter(r)
	_ = r.Run(":" + strconv.Itoa(core.DefaultHTTPPort))
}
