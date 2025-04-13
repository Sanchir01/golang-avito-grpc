package pvzgrpc

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	pvz "github.com/Sanchir01/golang-avito-proto/pkg/gen/golang/pvz"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServerAPI struct {
	pvz.UnimplementedPVZServiceServer
	primaryDB *pgxpool.Pool
}

func RegisterGrpcServer(gRPC *grpc.Server, db *pgxpool.Pool) {
	pvz.RegisterPVZServiceServer(gRPC, &ServerAPI{primaryDB: db})
}

func (s *ServerAPI) GetPVZList(ctx context.Context, req *pvz.GetPVZListRequest) (*pvz.GetPVZListResponse, error) {
	const op = "GetPVZList"

	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to acquire connection: %w", op, err)
	}
	defer conn.Release()

	query, args, err := sq.
		Select("id, registration_date, city").
		From("pvz").ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}
	defer rows.Close()

	var pvzlist []*pvz.PVZ

	for rows.Next() {
		var item pvz.PVZ
		var createdAt time.Time // Временная переменная для хранения времени

		// Сканируем в промежуточную переменную time.Time
		err := rows.Scan(&item.Id, &createdAt, &item.City)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}

		// Преобразуем time.Time в protobuf timestamp
		item.RegistrationDate = timestamppb.New(createdAt)

		pvzlist = append(pvzlist, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: error iterating rows: %w", op, err)
	}

	return &pvz.GetPVZListResponse{
		Pvzs: pvzlist,
	}, nil
}
