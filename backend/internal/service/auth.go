package service

import (
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/internal/storage"
	strgDto "avito-task-2025/backend/internal/storage/dto"

	// "avito-task-2025/backend/internal/storage/postgres"
	// "avito-task-2025/backend/pkg/jwt"
	"avito-task-2025/backend/pkg/jwt"
	"avito-task-2025/backend/pkg/logger"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrBadUsername       = errors.New("bad username")
	ErrBadPassword       = errors.New("bad password")
	ErrIncorrectPassword = errors.New("incorrect password")
)

var (
	ErrGenerateHashPass = errors.New("error while generate hash of the password")
	ErrInsertUserIntoDb = errors.New("error while insert neew user into db")
	ErrGenerateJWT      = errors.New("error while generate jwt token")
)

type AuthIntf interface {
	SignIn(ctx context.Context, request *svcDto.SignInRequest) (response *svcDto.SignInResponse, err error)
}

type AuthSvc struct {
	logger   logger.Interface
	userIntf storage.UserIntf
	jwtKey   string
}

func NewAuthSvc(logger logger.Interface, userIntf storage.UserIntf, jwtKey string) AuthIntf {
	return &AuthSvc{
		logger:   logger,
		userIntf: userIntf,
		jwtKey: jwtKey,
	}
}

func (s *AuthSvc) SignIn(ctx context.Context, request *svcDto.SignInRequest) (response *svcDto.SignInResponse, err error) {
	if request.Username == "" {
		s.logger.Warnf("%w", ErrBadUsername)
		return nil, ErrBadUsername
	}

	if request.Password == "" {
		s.logger.Warnf("%w", ErrBadPassword)
		return nil, ErrBadPassword
	}

	user, err := s.userIntf.GetByUsername(ctx, &strgDto.GetUserByUsernameRequest{
		Username: request.Username,
	})
	if err != nil {
		s.logger.Errorf(err.Error())
		return nil, err
	}

	if len(user.HashPassword) == 0 {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Errorf(ErrGenerateHashPass.Error())
			return nil, ErrGenerateHashPass
		}

		err = s.userIntf.Insert(ctx, &strgDto.InsertUserRequest{
			Username:     request.Username,
			HashPassword: hashPassword,
			CoinsAmount:  1000,
		})
		if err != nil {
			s.logger.Errorf(ErrInsertUserIntoDb.Error())
			return nil, ErrInsertUserIntoDb
		}
	} else {
		err = bcrypt.CompareHashAndPassword(user.HashPassword, []byte(request.Password))
		if err != nil {
			s.logger.Errorf(ErrIncorrectPassword.Error())
			return nil, ErrIncorrectPassword
		}
	}

	token, err := jwt.GenerateAuthToken(user.ID, s.jwtKey)
	if err != nil {
		s.logger.Errorf(ErrGenerateJWT.Error())
		return nil, ErrGenerateJWT
	}
	response = &svcDto.SignInResponse{
		JwtToken: token,
	}

	return response, nil

}
