package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sanchir01/golang-avito-grpc/internal/app"
)

func main() {
	env, err := app.NewEnv()
	if err != nil {
		panic(err)
	}
	env.Lg.Info("GRPC server started", "port", env.Cfg.Servers.Grpc.Port)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()
	go func() {
		env.GRPCSrv.MustRun()
	}()
	sign := <-ctx.Done()

	env.Lg.Info("received signal", "signal", sign)
	env.GRPCSrv.Stop()
	println("Hello World!")
}
