package service

import (
	"context"
	"github/eventApp/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userGetter interface {
	GetUserByUserName(userName string, ctx context.Context) (*models.User, error)
}

type LoginService struct {
	ug        userGetter
	jwtSecret string
}

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"-"`
}

func NewLoginService(jwtSecret string, ug userGetter) *LoginService {

	ls := &LoginService{
		ug:        ug,
		jwtSecret: jwtSecret,
	}

	return ls
}

func (ls *LoginService) Login(lr *LoginRequest, ctx context.Context) (*LoginResponse, error) {

	user, err := ls.ug.GetUserByUserName(lr.UserName, ctx)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(lr.Password))
	if err != nil {
		return nil, err
	}

	token, err := generateJWT(user.ID, ls.jwtSecret, time.Hour*24)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateJWT(user.ID, ls.jwtSecret, time.Hour*24*7)
	if err != nil {
		return nil, err
	}

	lresp := &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}

	return lresp, nil
}

func generateJWT(userID int64, jwtSecret string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     duration,
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
