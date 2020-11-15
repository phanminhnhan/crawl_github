package security

import (
	"github-trending/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const SecretKey = "fshjofjsdfo8oi3wyuf98wyu9876uhzxiou#@%%"

func Gentoken (user model.User)(string, error){
	claims := &model.JwtCustomClaims{
		UserId: user.UserId,
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour *24 ).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	resut, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err

	}
	return resut, nil
}