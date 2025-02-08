package services

import (
	"context"

	"github.com/CollabTED/CollabTed-Backend/pkg/types"
	"github.com/CollabTED/CollabTed-Backend/prisma"
	"github.com/CollabTED/CollabTed-Backend/prisma/db"
)

type LiveBoardService struct {
}

func NewLiveBoardService() *LiveBoardService {
	return &LiveBoardService{}
}

func (s *LiveBoardService) CreateBoard(data types.LiveBoardD) (*db.LiveBoardModel, error) {
	result, err := prisma.Client.LiveBoard.CreateOne(
		db.LiveBoard.Name.Set(data.Name),
		db.LiveBoard.Description.Set(data.Description),
		db.LiveBoard.Workspace.Link(
			db.Workspace.ID.Equals(data.WorkspaceId),
		),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	for _, user := range data.Users {
		usr, err := prisma.Client.UserWorkspace.FindFirst(
			db.UserWorkspace.UserID.Equals(user),
			db.UserWorkspace.WorkspaceID.Equals(data.WorkspaceId),
		).Exec(context.Background())
		if err != nil {
			return nil, err
		}
		_, err = prisma.Client.LiveBoard.FindUnique(
			db.LiveBoard.ID.Equals(result.ID),
		).Update(
			db.LiveBoard.Users.Link(
				db.UserWorkspace.ID.Equals(usr.ID),
			),
		).Exec(context.Background())

	}
	return result, nil
}

func (s *LiveBoardService) DeleteBoard(boardId string) (*db.LiveBoardModel, error) {
	result, err := prisma.Client.LiveBoard.FindUnique(
		db.LiveBoard.ID.Equals(boardId),
	).Delete().Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *LiveBoardService) GetBoard(boardId string) (*db.LiveBoardModel, error) {
	board, err := prisma.Client.LiveBoard.FindUnique(
		db.LiveBoard.ID.Equals(boardId),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return board, nil
}

func (s *LiveBoardService) GetWorkspaceBoards(workspaceId string) ([]db.LiveBoardModel, error) {
	board, err := prisma.Client.LiveBoard.FindMany(
		db.LiveBoard.WorkspaceID.Equals(workspaceId),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return board, nil
}
