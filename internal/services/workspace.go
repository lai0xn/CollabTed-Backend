package services

import (
	"context"

	"github.com/lai0xn/squid-tech/pkg/types"
	"github.com/lai0xn/squid-tech/prisma"
	"github.com/lai0xn/squid-tech/prisma/db"
)

type WorkspaceService struct{}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{}
}

func (s *WorkspaceService) CreateWorksapce(data types.WorkspaceD) (*db.WorkplaceModel, error) {
	result, err := prisma.Client.Workplace.CreateOne(
		db.Workplace.WorkplaceName.Set(data.Name),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *WorkspaceService) ListWorkspaces(userID string) ([]db.WorkplaceModel, error) {
	result, err := prisma.Client.Workplace.FindMany(
		db.Workplace.Users.Some(
			db.UserWorkplace.UserID.Equals(userID),
		),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}
