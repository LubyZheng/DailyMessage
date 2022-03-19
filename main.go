package main

import (
	"DailyMessage/web"
	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := env.Load()
	if err != nil {
		log.Fatalf("Load .env file error: %s", err)
	}
}

func main() {
	gin.DisableConsoleColor() //解决终端乱码
	router := web.SetupRouter()
	router.Run(":" + os.Getenv("SERVER_PORT"))
}
