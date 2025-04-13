package app

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net"

	pvzgrpc "github.com/Sanchir01/golang-avito-grpc/internal/grpc/pvz"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	log        *slog.Logger
	GRPCServer *grpc.Server
	port       int
}

func NewGRPCApp(
	log *slog.Logger,
	port int,
	db *pgxpool.Pool,
) *GRPCApp {
	GRPCServer := grpc.NewServer()
	pvzgrpc.RegisterGrpcServer(GRPCServer, db)
	return &GRPCApp{
		log:        log,
		GRPCServer: GRPCServer,
		port:       port,
	}
}

func (g *GRPCApp) MustRun() {
	if err := g.Run(); err != nil {
		panic(err)
	}
}

func (g *GRPCApp) Run() error {
	const op = "grpc.app.Run"
	log := g.log.With(slog.String("op", op), slog.Int("port", g.port))
	log.Info("Stating grpc")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	if err := g.GRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}

func (g *GRPCApp) Stop() error {
	const op = "grpc.app.Strop"
	g.log.With(slog.String("op", op), slog.Int("port", g.port))
	g.GRPCServer.GracefulStop()
	return nil
}
