package repository

import (
	"context"

	"githab.com/AlexeyBorisovets/ORDER/internal/model"
)

// type Repository interface {
// 	CreateUser(ctx context.Context, order *model.Order) (string, error)
// 	GetUserByID(ctx context.Context, userId string) (*model.User, error)
// 	GetAllUsers(ctx context.Context) ([]*model.User, error)
// 	DeleteUser(ctx context.Context, id string) error
// 	UpdateUser(ctx context.Context, id string, user *model.User) error
// 	UpdateAuth(ctx context.Context, id, refreshToken string) error
// 	SelectByIDAuth(ctx context.Context, id string) (model.User, error)
// 	GetUserByUserType(ctx context.Context, usertype string) ([]*model.User, error)
// 	GetBalanceByID(ctx context.Context, userId string) (uint, error)
// 	UpdateBalance(ctx context.Context, userId string, balance uint) error
// }
