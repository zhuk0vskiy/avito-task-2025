package service

import (
	svcDto "avito-task-2025/backend/internal/service/dto"
	"avito-task-2025/backend/internal/storage"
	strgDto "avito-task-2025/backend/internal/storage/dto"
	"time"

	"avito-task-2025/backend/pkg/jwt"
	"avito-task-2025/backend/pkg/logger"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyUsername     = errors.New("username is empty")
	ErrEmptyPassword     = errors.New("password is empty")
	ErrIncorrectPassword = errors.New("incorrect password")
)

var (
	ErrGenerateHashPass = errors.New("error while generate hash of the password")
	ErrGenerateJWT      = errors.New("error while generate jwt token")
)

type AuthIntf interface {
	SignIn(ctx context.Context, request *svcDto.SignInRequest) (response *svcDto.SignInResponse, err error)
}

type AuthSvc struct {
	logger     logger.Interface
	jwtManager jwt.ManagerIntf
	userIntf   storage.UserIntf
}

func NewAuthSvc(logger logger.Interface, jwtManager jwt.ManagerIntf, userIntf storage.UserIntf) AuthIntf {
	return &AuthSvc{
		logger:     logger,
		jwtManager: jwtManager,
		userIntf:   userIntf,
	}
}

func (s *AuthSvc) SignIn(ctx context.Context, request *svcDto.SignInRequest) (response *svcDto.SignInResponse, err error) {
	if request.Username == "" {
		s.logger.Warnf("%w", ErrEmptyUsername)
		return nil, ErrEmptyUsername
	}

	if request.Password == "" {
		s.logger.Warnf("%w", ErrEmptyPassword)
		return nil, ErrEmptyPassword
	}

	user, err := s.userIntf.GetByUsername(ctx, &strgDto.GetUserByUsernameRequest{
		Username: request.Username,
	})
	if err != nil {
		s.logger.Warnf(err.Error())
		return nil, err
	}

	if len(user.HashPassword) == 0 {
		s.logger.Debugf("before generate %s", time.Now().UnixMilli())
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
		if err != nil {
			s.logger.Errorf(ErrGenerateHashPass.Error())
			return nil, ErrGenerateHashPass
		}
		s.logger.Debugf("after generate %s", time.Now().UnixMilli())

		response, err := s.userIntf.Insert(ctx, &strgDto.InsertUserRequest{
			Username:     request.Username,
			HashPassword: hashPassword,
		})
		if err != nil {
			s.logger.Warnf(err.Error())
			return nil, err
		}
		user.ID = response.ID
	} else {
		s.logger.Debugf("before compare %s", time.Now().UnixMilli())
		err = bcrypt.CompareHashAndPassword(user.HashPassword, []byte(request.Password))
		if err != nil {
			s.logger.Warnf(ErrIncorrectPassword.Error())
			return nil, ErrIncorrectPassword
		}
		s.logger.Debugf("after compare %s", time.Now().UnixMilli())
	}

	token, err := s.jwtManager.GenerateAuthToken(user.ID)
	if err != nil {
		s.logger.Errorf(ErrGenerateJWT.Error())
		return nil, ErrGenerateJWT
	}

	response = &svcDto.SignInResponse{
		JwtToken: token,
	}

	return response, nil

}
