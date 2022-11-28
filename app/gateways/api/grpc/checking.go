package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	pb "github.com/jpgsaraceni/suricate-bank/protos/checking"
)

type checkingServer struct {
	pb.UnimplementedCheckingServer
	uc accountuc.Usecase
}

func (s *checkingServer) GetBalance(ctx context.Context, acc *pb.Account) (*pb.Balance, error) {
	UUID, err := uuid.Parse(acc.Id)
	if err != nil {
		return nil, err
	}

	balance, err := s.uc.GetBalance(ctx, account.ID(UUID))
	if err != nil {
		return nil, err
	}

	return &pb.Balance{Amount: int32(balance)}, nil
}
