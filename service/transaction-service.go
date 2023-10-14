package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/repository"
	"github.com/midtrans/midtrans-go"

	"github.com/google/uuid"
	"github.com/mashingan/smapping"
)

type DepositService interface {
	InsertDeposit(Deposit dto.DepositDTO) entity.Deposit
	All(page int, pageSize int) ([]entity.Deposit, error)
	FindDepositByIDUser(idUiser uint64, int, pageSize int) ([]entity.Deposit, error)
	FindDepositByID(id uint64) *entity.Deposit
	SaveFile(file *multipart.FileHeader) (string, error)
	TotalDeposit() int64
	TotalDepositByUserID(idUser uint64) int64
	InsertPaymentToken(transactionID uint64, paymentToken string, virtualAcc string) error
	UpdateDepositStatus(orderID uint64, newStatus uint64) error
}

type transactionService struct {
	DepositRepository repository.DepositRepository
}

func NewDepositService(fundRep repository.DepositRepository) DepositService {
	return &transactionService{
		DepositRepository: fundRep,
	}
}

func (service *transactionService) InsertDeposit(b dto.DepositDTO) entity.Deposit {
	Deposit := entity.Deposit{}
	err := smapping.FillStruct(&Deposit, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	res := service.DepositRepository.InsertDeposit(&Deposit)
	return res
}

func (service *transactionService) InsertPaymentToken(transactionID uint64, paymentToken string, virtualAcc string) error {
	err := service.DepositRepository.StorePaymentToken(transactionID, paymentToken, virtualAcc)
	if err != nil {
		if midErr, ok := err.(*midtrans.Error); ok {
			return errors.New(midErr.Message)
		}
		return err
	}
	return nil
}

func (service *transactionService) TotalDeposit() int64 {
	return service.DepositRepository.TotalDeposit()
}

func (service *transactionService) TotalDepositByUserID(idUser uint64) int64 {
	return service.DepositRepository.TotalDepositByUserID(idUser)
}

func (service *transactionService) All(page int, pageSize int) ([]entity.Deposit, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	return service.DepositRepository.All(page, pageSize)
}

func (service *transactionService) FindDepositByIDUser(idUser uint64, page int, pageSize int) ([]entity.Deposit, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	return service.DepositRepository.FindDepositByIDUser(idUser, page, pageSize)
}

func (service *transactionService) FindDepositByID(id uint64) *entity.Deposit {
	return service.DepositRepository.FindDepositByID(id)
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

func (service *transactionService) UpdateDepositStatus(orderID uint64, newStatus uint64) error {
	// Fetch the MasterJual entity by order ID
	masterJual := service.DepositRepository.FindDepositByID(orderID)
	if masterJual.ID == 0 {
		return fmt.Errorf("MasterJual not found for order ID %s", orderID)
	}

	masterJual.Status = newStatus

	err := service.DepositRepository.UpdateDepositStatus(masterJual.ID, newStatus)
	if err != nil {
		return err
	}

	return nil
}
