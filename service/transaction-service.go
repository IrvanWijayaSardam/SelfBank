package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/repository"

	"github.com/mashingan/smapping"
)

type TransactionService interface {
	InsertTransaction(Transaction dto.TransactionDTO) entity.Transaction
	All(page int, pageSize int) ([]entity.Transaction, error)
	FindTransactionByIDUser(idUiser uint64, int, pageSize int) ([]entity.Transaction, error)
	FindTransactionByID(id uint64) *entity.Transaction
	TotalTransaction() int64
	TotalTransactionByUserID(idUser uint64) int64
	UpdateTransactionStatus(orderID uint64, newStatus uint64) error
	ValidateAccNumber(accNumber uint64) bool
}

type transactionService struct {
	TransactionRepository repository.TransactionRepository
}

func NewTransactionService(fundRep repository.TransactionRepository) TransactionService {
	return &transactionService{
		TransactionRepository: fundRep,
	}
}

func (service *transactionService) InsertTransaction(b dto.TransactionDTO) entity.Transaction {
	Transaction := entity.Transaction{}
	err := smapping.FillStruct(&Transaction, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	res := service.TransactionRepository.InsertTransaction(&Transaction)
	return res
}

func (service *transactionService) TotalTransaction() int64 {
	return service.TransactionRepository.TotalTransaction()
}

func (service *transactionService) ValidateAccNumber(accNumber uint64) bool {
	return service.TransactionRepository.ValidateAccNumber(accNumber)
}

func (service *transactionService) TotalTransactionByUserID(idUser uint64) int64 {
	return service.TransactionRepository.TotalTransactionByUserID(idUser)
}

func (service *transactionService) All(page int, pageSize int) ([]entity.Transaction, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	return service.TransactionRepository.All(page, pageSize)
}

func (service *transactionService) FindTransactionByIDUser(idUser uint64, page int, pageSize int) ([]entity.Transaction, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	return service.TransactionRepository.FindTransactionByIDUser(idUser, page, pageSize)
}

func (service *transactionService) FindTransactionByID(id uint64) *entity.Transaction {
	return service.TransactionRepository.FindTransactionByID(id)
}

func (service *transactionService) UpdateTransactionStatus(orderID uint64, newStatus uint64) error {
	// Fetch the MasterJual entity by order ID
	masterJual := service.TransactionRepository.FindTransactionByID(orderID)
	if masterJual.ID == 0 {
		return fmt.Errorf("MasterJual not found for order ID %s", orderID)
	}

	masterJual.Status = newStatus

	err := service.TransactionRepository.UpdateTransactionStatus(masterJual.ID, newStatus)
	if err != nil {
		return err
	}

	return nil
}
