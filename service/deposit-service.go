package service

import (
	"bytes"
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
	"github.com/jung-kurt/gofpdf"
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
	GenerateDepositPDF(deposits []dto.DepositResponse) (*bytes.Buffer, error)
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

func (service *depositService) GenerateDepositPDF(Deposits []dto.DepositResponse) (*bytes.Buffer, error) {
	pdf := gofpdf.New("L", "mm", "A2", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Deposit Report")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 10, "Deposit ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "ID User", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Date", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 10, "Amount", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 10, "Status", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 10, "Virtual Account", "1", 0, "C", false, 0, "")
	pdf.CellFormat(220, 10, "URL Callback", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	for _, deposit := range Deposits {
		pdf.CellFormat(40, 10, deposit.Id_deposit, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, helper.Uint64ToString(deposit.Id_user), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, deposit.Date.Format("2006-01-02 15:04:05"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 10, helper.Uint64ToString(deposit.Amount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 10, deposit.Status, "1", 0, "C", false, 0, "")
		pdf.CellFormat(50, 10, deposit.Virtual_account, "1", 0, "C", false, 0, "")
		pdf.CellFormat(220, 10, deposit.Url_callback, "1", 1, "C", false, 0, "")
	}

	pdfBuffer := new(bytes.Buffer)

	err := pdf.Output(pdfBuffer)
	if err != nil {
		return nil, err
	}

	return pdfBuffer, nil
}
