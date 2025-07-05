package main

import (
	"log"
	"net/http"
	"os"

	"github.com/asm2212/gochat/internal"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := internal.NewRedisClient(redisAddr)
	userSvc := internal.NewUserService(rdb)
	chatSvc := internal.NewChatService(rdb)

	handler := internal.NewHandler(userSvc, chatSvc)
	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
