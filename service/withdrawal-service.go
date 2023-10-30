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

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/helper"
	"github.com/IrvanWijayaSardam/SelfBank/repository"
	"github.com/jung-kurt/gofpdf"

	"github.com/google/uuid"
	"github.com/mashingan/smapping"
)

type WithdrawalService interface {
	InsertWithdrawal(Withdrawal dto.WithdrawalDTO) entity.Withdrawal
	All(page int, pageSize int) ([]entity.Withdrawal, error)
	FindWithdrawalByIDUser(idUiser uint64, int, pageSize int) ([]entity.Withdrawal, error)
	FindWithdrawalByID(id uint64) *entity.Withdrawal
	SaveFile(file *multipart.FileHeader) (string, error)
	TotalWithdrawal() int64
	TotalWithdrawalByUserID(idUser uint64) int64
	UpdateWithdrawalStatus(orderID uint64, newStatus uint64) error
	GenerateWithdrawalPDF(Transactions []entity.Withdrawal) (*bytes.Buffer, error)
}

type withdrawalService struct {
	WithdrawalRepository repository.WithdrawalRepository
}

func NewWithdrawalService(fundRep repository.WithdrawalRepository) WithdrawalService {
	return &withdrawalService{
		WithdrawalRepository: fundRep,
	}
}

func (service *withdrawalService) InsertWithdrawal(b dto.WithdrawalDTO) entity.Withdrawal {
	Withdrawal := entity.Withdrawal{}
	err := smapping.FillStruct(&Withdrawal, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	res := service.WithdrawalRepository.InsertWithdrawal(&Withdrawal)
	return res
}

func (service *withdrawalService) TotalWithdrawal() int64 {
	return service.WithdrawalRepository.TotalWithdrawal()
}

func (service *withdrawalService) TotalWithdrawalByUserID(idUser uint64) int64 {
	return service.WithdrawalRepository.TotalWithdrawalByUserID(idUser)
}

func (service *withdrawalService) All(page int, pageSize int) ([]entity.Withdrawal, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	return service.WithdrawalRepository.All(page, pageSize)
}

func (service *withdrawalService) FindWithdrawalByIDUser(idUser uint64, page int, pageSize int) ([]entity.Withdrawal, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	return service.WithdrawalRepository.FindWithdrawalByIDUser(idUser, page, pageSize)
}

func (service *withdrawalService) FindWithdrawalByID(id uint64) *entity.Withdrawal {
	return service.WithdrawalRepository.FindWithdrawalByID(id)
}

func (service *withdrawalService) SaveFile(file *multipart.FileHeader) (string, error) {
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

func (service *withdrawalService) UpdateWithdrawalStatus(orderID uint64, newStatus uint64) error {
	// Fetch the MasterJual entity by order ID
	masterJual := service.WithdrawalRepository.FindWithdrawalByID(orderID)
	if masterJual.ID == 0 {
		return fmt.Errorf("MasterJual not found for order ID %s", orderID)
	}

	masterJual.Status = newStatus

	err := service.WithdrawalRepository.UpdateWithdrawalStatus(masterJual.ID, newStatus)
	if err != nil {
		return err
	}

	return nil
}

func (service *withdrawalService) GenerateWithdrawalPDF(Transactions []entity.Withdrawal) (*bytes.Buffer, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Withdrawal Report")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 10, "Withdrawal ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "ID User", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 10, "Date", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Amount", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 10, "Transaction To", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	for _, transaction := range Transactions {
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", transaction.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", transaction.ID_User), "1", 0, "C", false, 0, "")
		pdf.CellFormat(50, 10, helper.ConvertUnixtime(transaction.Date).Format("2006-01-02 15:04:05"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, helper.Uint64ToString(transaction.Amount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 10, transaction.To, "1", 1, "C", false, 0, "")
	}

	pdfBuffer := new(bytes.Buffer)

	err := pdf.Output(pdfBuffer)
	if err != nil {
		return nil, err
	}

	return pdfBuffer, nil
}
