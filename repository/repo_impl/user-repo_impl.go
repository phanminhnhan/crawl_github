package repo_impl

import (
	"database/sql"
	"github-trending/database"
	"github-trending/log"
	"github-trending/banana"
	"github-trending/model"
	"github-trending/model/request"
	"github-trending/repository"
	"context"
	"github.com/lib/pq"

	"time"
)

type UserRepoImpl struct{
	sql *database.Sql

}


func NewRepo(sql *database.Sql) repository.UserRepo{
	 return &UserRepoImpl{
	 	sql: sql,
	 }
}


func (u *UserRepoImpl) AddUser(ctx context.Context, user model.User)(model.User, error){
	queryStatement := `INSERT INTO users(user_id,  full_name , email, password, role,  created_at, updated_at )
	VALUES (:user_id,  :full_name , :email, :password, :role,  :created_at, :updated_at)	
	
`
	user.CreateAt = time.Now()
	user.UpdateAt = time.Now()
	_ , err := u.sql.Db.NamedExecContext(ctx , queryStatement, user)
	if err != nil{
		log.Error(err.Error())
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name()== "unique_violation" {
				return user, banana.UserConflict
			}
		}
		return user , banana.SignUpFail
	}
	return user, nil
}

func (u *UserRepoImpl)Checklogin(ctx context.Context, loginReq request.ReqLogin)(model.User, error){
	user := model.User{}
	queryStatement := "SELECT * FROM users WHERE email = $1"
	err := u.sql.Db.GetContext(ctx ,&user,  queryStatement, loginReq.Email)
	if err != nil {

		if err == sql.ErrNoRows{
			return model.User{}, banana.UserNotFound
		}
		log.Error(err.Error())
		return user, err
	}
	return model.User{}, nil

}

func (u *UserRepoImpl)GetUerByID(ctx context.Context, userId string )(model.User, error){
	user := model.User{}
	err := u.sql.Db.GetContext(ctx, &user, "SELECT * FROM users WHERE user_id= $1", userId)
	if err != nil{
		if err == sql.ErrNoRows{
			return user, banana.UserNotFound
		}
		log.Error(err.Error())
		return user, err
	}
	return user, nil
}

func(u *UserRepoImpl)UpdateUser(ctx context.Context, user model.User)(model.User, error){
	sqlStatement := `
		UPDATE users
		SET 
			full_name  = (CASE WHEN LENGTH(:full_name) = 0 THEN full_name ELSE :full_name END),
			email = (CASE WHEN LENGTH(:email) = 0 THEN email ELSE :email END),
			updated_at 	  = COALESCE (:updated_at, updated_at)
		WHERE user_id    = :user_id
	`

	user.UpdateAt= time.Now()
	result, err := u.sql.Db.NamedExecContext(ctx, sqlStatement, user)
	if err != nil {
		log.Error(err.Error())
		return user, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return user, banana.UserNotUpdated
	}
	if count == 0 {
		return user, banana.UserNotUpdated
	}

	return user, nil
}