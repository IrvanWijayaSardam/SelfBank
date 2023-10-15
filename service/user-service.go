package service

import (
	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/repository"
)

type UserService interface {
	FindUser(id uint64) entity.User
	GetSaldo(idUser uint64) int64
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRep repository.UserRepository) UserService {
	return &userService{
		userRepository: userRep,
	}
}

func (service *userService) FindUser(id uint64) entity.User {
	return service.userRepository.ProfileUser(id)
}

func (service *userService) GetSaldo(id uint64) int64 {
	totalDeposit := service.userRepository.TotalDepositByUserID(id)
	totalWithdrawal := service.userRepository.TotalWithdrawalByUserID(id)
	balance := totalDeposit - totalWithdrawal

	return balance
}
