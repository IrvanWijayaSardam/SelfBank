package service

import (
	"strconv"

	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/repository"
)

type UserService interface {
	All(page int, pageSize int) ([]entity.User, error)
	FindUser(id uint64) entity.User
	GetSaldo(idUser uint64) int64
	UpdateUser(user entity.User) entity.User
	DeleteUser(idUser uint64) bool
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRep repository.UserRepository) UserService {
	return &userService{
		userRepository: userRep,
	}
}

func (service *userService) All(page int, pageSize int) ([]entity.User, error) {
	return service.userRepository.All(page, pageSize)
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
	totalTransactionIn := service.userRepository.TotalTransactionInByAccountNumber(strconv.Itoa(int(user.AccountNumber)))
	totalTransactionOut := service.userRepository.TotalTransactionFromByAccountNumber(strconv.Itoa(int(user.AccountNumber)))
	balance := (totalDeposit + totalTransactionIn) - totalWithdrawal - totalTransactionOut

	return balance
}

func (service *userService) DeleteUser(id uint64) bool {
	return service.userRepository.DeleteUser(id)
}
