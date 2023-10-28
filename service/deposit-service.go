package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/helper"
	"github.com/IrvanWijayaSardam/SelfBank/repository"
	"github.com/midtrans/midtrans-go"

	"github.com/google/uuid"
	"github.com/mashingan/smapping"
)

type DepositService interface {
	InsertDeposit(Deposit dto.DepositDTO) entity.Deposit
	All(page int, pageSize int) ([]entity.Deposit, error)
	FindDepositByIDUser(idUser uint64, int, pageSize int) ([]entity.Deposit, error)
	FindDepositByID(id string) *entity.Deposit
	SaveFile(file *multipart.FileHeader) (string, error)
	TotalDeposit() int64
	TotalDepositByUserID(idUser uint64) int64
	InsertPaymentToken(transactionID string, paymentToken string, virtualAcc string, callbackUrl string) error
	UpdateDepositStatus(orderID string, newStatus uint64) error
	FindPaymentInfoById(depositId string) *entity.PaymentToken
}

type depositService struct {
	DepositRepository repository.DepositRepository
}

func NewDepositService(fundRep repository.DepositRepository) DepositService {
	return &depositService{
		DepositRepository: fundRep,
	}
}

func (service *depositService) InsertDeposit(b dto.DepositDTO) entity.Deposit {
	Deposit := entity.Deposit{}
	err := smapping.FillStruct(&Deposit, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	numInt := int(helper.GenerateTrxId())
	Deposit.ID = strconv.Itoa(numInt)
	res := service.DepositRepository.InsertDeposit(&Deposit)
	return res
}

func (service *depositService) InsertPaymentToken(transactionID string, paymentToken string, virtualAcc string, callbackUrl string) error {
	err := service.DepositRepository.StorePaymentToken(transactionID, paymentToken, virtualAcc, callbackUrl)
	if err != nil {
		if midErr, ok := err.(*midtrans.Error); ok {
			return errors.New(midErr.Message)
		}
		return err
	}
	return nil
}

func (service *depositService) TotalDeposit() int64 {
	return service.DepositRepository.TotalDeposit()
}

func (service *depositService) TotalDepositByUserID(idUser uint64) int64 {
	return service.DepositRepository.TotalDepositByUserID(idUser)
}

func (service *depositService) All(page int, pageSize int) ([]entity.Deposit, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	return service.DepositRepository.All(page, pageSize)
}

func (service *depositService) FindDepositByIDUser(idUser uint64, page int, pageSize int) ([]entity.Deposit, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	return service.DepositRepository.FindDepositByIDUser(idUser, page, pageSize)
}

func (service *depositService) FindDepositByID(id string) *entity.Deposit {
	return service.DepositRepository.FindDepositByID(id)
}

func (service *depositService) FindPaymentInfoById(id string) *entity.PaymentToken {
	return service.DepositRepository.FindPaymentInfoById(id)
}

func (service *depositService) SaveFile(file *multipart.FileHeader) (string, error) {
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

func (service *depositService) UpdateDepositStatus(orderID string, newStatus uint64) error {
	// Fetch the MasterJual entity by order ID
	masterJual := service.DepositRepository.FindDepositByID(orderID)
	if masterJual.ID == "0" {
		return fmt.Errorf("MasterJual not found for order ID %s", orderID)
	}

	masterJual.Status = newStatus

	err := service.DepositRepository.UpdateDepositStatus(masterJual.ID, newStatus)
	if err != nil {
		return err
	}

	return nil
}
