package service

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/repository"

	"github.com/google/uuid"
	"github.com/mashingan/smapping"
)

type TransactionService interface {
	InsertTransaction(Transaction dto.TransactionDTO) entity.Transaction
	All(page int, pageSize int) ([]entity.Transaction, error)
	FindTransactionByIDUser(idUiser uint64, int, pageSize int) ([]entity.Transaction, error)
	FindTransactionByID(id string) *entity.Transaction
	SaveFile(file *multipart.FileHeader) (string, error)
	TotalTransaction() int64
	TotalTransactionByUserID(idUser uint64) int64
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

func (service *transactionService) FindTransactionByID(id string) *entity.Transaction {
	return service.TransactionRepository.FindTransactionByID(id)
}

func (service *transactionService) SaveFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	cdnDir := "cdn" // Specify the desired directory name

	// Create the cdn directory if it doesn't exist in the current working directory
	err = os.MkdirAll(cdnDir, 0755)
	if err != nil {
		return "", err
	}

	// Generate a random file name using UUID and append the original file extension
	fileExt := filepath.Ext(file.Filename)
	fileName := uuid.New().String() + fileExt
	filePath := filepath.Join(cdnDir, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
