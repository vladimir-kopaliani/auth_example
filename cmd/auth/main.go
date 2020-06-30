package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	serverhttp "github.com/vladimir-kopaliani/auth_example/internal/http-server"
	authrepo "github.com/vladimir-kopaliani/auth_example/internal/repository"
	"github.com/vladimir-kopaliani/auth_example/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// handle interupt signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		select {
		case <-signalChan:
			log.Println("Got Interrupt signal. Shutting down...")
			cancel()
		}
	}()

	// repository
	repo, err := authrepo.New(ctx, &authrepo.Configuration{
		URI: os.Getenv("DB_ADDRESS"),
	})
	if err != nil {
		log.Println(err)
		cancel()
	}

	// service
	serv, err := service.New(ctx, service.Configuration{
		Repository: repo,
	})
	if err != nil {
		log.Println(err)
		cancel()
	}

	// http server
	httpServer, err := serverhttp.New(
		ctx,
		serverhttp.Configuration{
			Port:    os.Getenv("HTTP_PORT"),
			Service: serv,
		})
	if err != nil {
		log.Println(err)
		cancel()
	}
	go func() {
		err = httpServer.Launch(ctx)
		if err != nil {
			log.Println(err)
			cancel()
		}
	}()
	defer httpServer.Close(ctx)

	<-ctx.Done()
}
