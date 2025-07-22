package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/johnjiangtw0804/pitngo-rate-limit/env"
	"github.com/johnjiangtw0804/pitngo-rate-limit/infra"
	router "github.com/johnjiangtw0804/pitngo-rate-limit/router"
	"github.com/spf13/viper"
)

func main() {
	var err error
	loc, err := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	if err != nil {
		log.Fatalf("Invalid timezone: %v", err)
	}
	time.Local = loc

	env, err := env.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig fail: %v", err)
	}
	// log.Default().Println(env)
	redisClient, err := infra.ConnectDBS(env)
	if err != nil {
		// %v is a general-purpose formatting verb used
		log.Fatalf("Connect Redis fail: %v", err)
	}
	// log.Default().Println(redisClient)

	// test with simple set and get and unlink Redis command
	ctx := context.Background()
	err = redisClient.Set(ctx, "Jonathan", "He is Hot", 0).Err()
	if err != nil {
		log.Fatalf("Redis SET command failed: %v", err)
	}
	log.Default().Printf("key: %s", "Jonathan")

	val, err := redisClient.Get(ctx, "Jonathan").Result()
	if err != nil {
		log.Fatalf("Redis GET command failed: %v", err)
	}

	log.Default().Printf("val: %s", val)

	if env.AppEnv == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router, err := router.RegisterRouter(env, redisClient)
	if err != nil {
		log.Fatalf("Router setup failed: %v", err)
	}

	srv := &http.Server{
		Addr:    env.AppPort,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal, 1)

	// relay incoming signals to quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down Server ...")

	// a timeout of 5 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	} else {
		log.Println("Server exited properly")
	}
}
