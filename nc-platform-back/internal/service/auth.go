package service

import (
	"github.com/go-chi/jwtauth"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"nc-platform-back/internal/config"
	"nc-platform-back/internal/repository"
)

type AuthService struct {
	logger         logrus.FieldLogger
	config         *config.Config
	userRepository *repository.UserRepository
}

func NewAuthService(config *config.Config, userRepository *repository.UserRepository) *AuthService {
	return &AuthService{
		logger:         logrus.New(),
		config:         config,
		userRepository: userRepository,
	}
}

func (s *AuthService) Login(username string, password string) (string, error) {
	user, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		s.logger.Error(err)
		return "", err
	}

	s.logger.Info(user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	logrus.Info(err)
	if err != nil {
		return "", err
	}

	jwtAuth := jwtauth.New(s.config.JwtAlgo, []byte(s.config.JwtKey), nil)
	_, jwt, err := jwtAuth.Encode(map[string]interface{}{"userId": user.Id})
	if err != nil {
		s.logger.Error(err)
		return "", err
	}
	return jwt, nil
}
