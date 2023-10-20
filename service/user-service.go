package service

import (
	"strconv"

	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/repository"
)

type UserService interface {
	FindUser(id uint64) entity.User
	GetSaldo(idUser uint64) int64
	UpdateUser(user entity.User) entity.User
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

func (service *userService) UpdateUser(user entity.User) entity.User {
	return service.userRepository.UpdateUser(user)
}

func (service *userService) GetSaldo(id uint64) int64 {
	user := service.userRepository.ProfileUser(id)

	totalDeposit := service.userRepository.TotalDepositByUserID(id)
	totalWithdrawal := service.userRepository.TotalWithdrawalByUserID(id)
	totalTransactionIn := service.userRepository.TotalTransactionByAccountNumber(strconv.Itoa(int(user.AccountNumber)))
	balance := (totalDeposit + totalTransactionIn) - totalWithdrawal

	return balance
}
