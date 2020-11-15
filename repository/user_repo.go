package repository

import (
	"context"
	"github-trending/model"
	"github-trending/model/request"
)

type UserRepo interface {
	AddUser(ctx context.Context, user model.User)(model.User, error)
	Checklogin(ctx context.Context, loginReq request.ReqLogin)(model.User, error)
	GetUerByID(ctx context.Context, userId string )(model.User, error)
	UpdateUser(ctx context.Context, user model.User)(model.User, error)
}