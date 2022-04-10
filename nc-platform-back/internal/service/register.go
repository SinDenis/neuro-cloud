package service

import (
	"demo-rest/internal/domain"
	"demo-rest/internal/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type RegisterService struct {
	logger         logrus.FieldLogger
	userRepository *repository.UserRepository
}

func NewRegisterService(userRepository *repository.UserRepository) *RegisterService {
	return &RegisterService{
		logger:         logrus.New(),
		userRepository: userRepository,
	}
}

func (s *RegisterService) Register(registerUser domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Error hash password", err)
		return err
	}

	registerUser.Password = string(hashedPassword)
	err = s.userRepository.SaveUser(registerUser)
	if err != nil {
		s.logger.Error("Failed save new user in db user", err)
		return err
	}
	return nil
}
