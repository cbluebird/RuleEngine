package main

import (
	"engine/app/midwares"
	"engine/config/config"
	"engine/config/router"
	"engine/config/session"
	"engine/config/startConfig"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	startConfig.Init()
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(midwares.ErrHandler())
	r.NoMethod(midwares.HandleNotFound)
	r.NoRoute(midwares.HandleNotFound)
	session.Init(r)
	router.Init(r)
	err := r.Run(":" + config.Config.GetString("server.port"))
	if err != nil {
		log.Fatal("ServerStartFailed", err)
	}
}
