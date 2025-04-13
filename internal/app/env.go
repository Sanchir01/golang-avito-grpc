package app

import (
	"context"
	"log/slog"

	"github.com/Sanchir01/golang-avito-grpc/internal/config"
)

type Env struct {
	Lg       *slog.Logger
	Cfg      *config.Config
	Database *Database
	GRPCSrv  *GRPCApp
}

func NewEnv() (*Env, error) {
	cfg := config.MustLoadConfig()
	logger := setupLogger(cfg.Env)
	ctx := context.Background()

	primarydb, err := NewDataBases(ctx, cfg.PrimaryDB.User, cfg.PrimaryDB.Host, cfg.PrimaryDB.Dbname, cfg.PrimaryDB.Port, cfg.PrimaryDB.MaxAttempts)
	if err != nil {
		return nil, err
	}
	grpcApp := NewGRPCApp(logger, cfg.Servers.Grpc.Port, primarydb.PrimaryDB)
	env := Env{
		Lg:       logger,
		Cfg:      cfg,
		Database: primarydb,
		GRPCSrv:  grpcApp,
	}

	return &env, nil
}
