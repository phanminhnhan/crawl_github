package handler

import (
	"github-trending/log"
	"github-trending/model"
	"github-trending/model/request"
	"github-trending/repository"
	"github-trending/security"
	"github.com/dgrijalva/jwt-go"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"net/http"
)


type UserHandler struct {
	UserRepo repository.UserRepo
}


func (u *UserHandler)HandleSignUp(c echo.Context) error {
	req := new(request.ReqRes)
	err := c.Bind(&req)
	if err != nil{
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message: err.Error(),
			Data: nil,
		})
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil{
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message: "nhap thieu thong tin",
			Data: nil,
		})
	}
	//hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	//if err != nil {
	//	log.Error(err.Error())
	//}
	//hash , _ := security.HashPassword(req.Password)
	isEmail := setting.IsEmailValid(req.Email)
	if !isEmail {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message: "Email is in valid",
			Data: nil,
		})
	}
	hash := security.HashAndSalt([]byte(req.Password))
	role := model.MEMBER.String()
	userId, err := uuid.NewUUID()
	if err != nil{
		log.Error(err.Error())
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message: err.Error(),
			Data: nil,
		})
	}

	
	user := model.User{
		UserId:   userId.String(),
		FullName: req.Fullname,
		Email:    req.Email,
		Password: hash,
		Role:     role,
		Token:    "",
	}
	result, err := u.UserRepo.AddUser(c.Request().Context(), user )
	if err != nil{
		log.Error(err.Error())
		return c.JSON(http.StatusConflict, model.Response{
			StatusCode: http.StatusConflict,
			Message: err.Error(),
			Data: nil,
		})
	}

	token, err := security.Gentoken(user)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message: err.Error(),
			Data: nil,
		})

	}
	result.Token = token

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message: "done",
		Data: result,
	})
}
func (u *UserHandler) HandleSignIn(c echo.Context) error {
	req := request.ReqLogin{}
	err := c.Bind(&req)
	if err != nil{
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message: "loi1",//err.Error(),
			Data: nil,
		})
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil{
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message: "loi2",//err.Error(),
			Data: nil,
		})
	}
	user, err := u.UserRepo.Checklogin(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message: "loi3",//err.Error(),
			Data: nil,
		})
	}

	// check password trong data base

	isTheSame := security.ComparePasswords(user.Password, []byte(req.Password))
	if !isTheSame {
		log.Error()
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Đăng nhập thất bại",
			Data:       nil,
		})
	}

	token, err := security.Gentoken(user)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message: err.Error(),
			Data: nil,
		})

	}
	user.Token = token
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message: "Log in Done",
		Data: user,
	})
}


func (u *UserHandler)TestJWT(c echo.Context)error{
	return c.String(http.StatusOK, "test oke")
}

func (u *UserHandler)UserProfile(c echo.Context)error{
	tokenData := c.Get("user").(*jwt.Token)
	claims := tokenData.Claims.(*model.JwtCustomClaims)
	userID := claims.UserId
	user, err := u.UserRepo.GetUerByID(c.Request().Context(), userID)
	if err != nil{
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message: err.Error(),
			Data: nil,
		})
	}
	return c.JSON(http.StatusInternalServerError, model.Response{
		StatusCode: 200,
		Message: "done",
		Data: user,
	})

}


func (u *UserHandler)UpdateUser(c echo.Context, )error{

	req := request.ReqUpdate{}
	err := c.Bind(&req)
	if err != nil{
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message: "loi1",//err.Error(),
			Data: nil,
		})
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil{
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message: "nhap thieu thong tin",
			Data: nil,
		})
	}
	tokenData := c.Get("user").(*jwt.Token)
	claims := tokenData.Claims.(*model.JwtCustomClaims)
	user := model.User{
		UserId: claims.UserId,
		FullName: req.Fullname,
		Email: req.Email,
	}
	user, err = u.UserRepo.UpdateUser(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, model.Response{
		StatusCode: http.StatusCreated,
		Message:    "done",
		Data:       user,
	})

}