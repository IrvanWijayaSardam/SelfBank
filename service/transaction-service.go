package service

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/helper"
	"github.com/IrvanWijayaSardam/SelfBank/repository"
	"github.com/jung-kurt/gofpdf"

	"github.com/mashingan/smapping"
)

type TransactionService interface {
	InsertTransaction(Transaction dto.TransactionDTO) entity.Transaction
	All(page int, pageSize int) ([]entity.Transaction, error)
	FindTransactionByIDUser(idUiser uint64, int, pageSize int) ([]entity.Transaction, error)
	FindTransactionByID(id uint64) entity.Transaction
	TotalTransaction() int64
	TotalTransactionByUserID(idUser uint64) int64
	UpdateTransactionStatus(orderID uint64, newStatus uint64) error
	ValidateAccNumber(accNumber uint64) bool
	GenerateTransactionPDF(Transactions []entity.Transaction) (*bytes.Buffer, error)
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

func (service *transactionService) FindTransactionByID(id uint64) entity.Transaction {
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

func (service *transactionService) GenerateTransactionPDF(Transactions []entity.Transaction) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Transaction Report")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 10, "Transaction ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Amount", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 10, "Date", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 10, "Transaction To", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	for _, transaction := range Transactions {
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", transaction.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, helper.Uint64ToString(transaction.Amount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(50, 10, helper.ConvertUnixtime(transaction.Date).Format("2006-01-02 15:04:05"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 10, helper.Uint64ToString(transaction.TransactionTo), "1", 1, "C", false, 0, "")
	}

	pdfBuffer := new(bytes.Buffer)

	err := pdf.Output(pdfBuffer)
	if err != nil {
		return nil, err
	}

	return pdfBuffer, nil
}
