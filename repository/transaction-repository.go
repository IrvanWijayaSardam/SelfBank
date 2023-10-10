package repository

import (
	"errors"

	"github.com/IrvanWijayaSardam/SelfBank/entity"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	InsertTransaction(brg *entity.Transaction) entity.Transaction
	All(page int, pageSize int) ([]entity.Transaction, error)
	UpdateTransaction(plg entity.Transaction) entity.Transaction
	FindTransactionByID(id uint64) *entity.Transaction
	FindTransactionByIDUser(id uint64, page int, pageSize int) ([]entity.Transaction, error)
	TotalTransaction() int64
	TotalTransactionByUserID(idUser uint64) int64
	StorePaymentToken(transactionID uint64, paymentToken string, virtualAcc string) error
	UpdateTransactionStatus(id uint64, newStatus uint64) error
}

type TransactionConnection struct {
	connection *gorm.DB
}

func (db *TransactionConnection) TotalTransaction() int64 {
	var count int64
	result := db.connection.Model(&entity.Transaction{}).Where("status = ?", 1).Count(&count)
	if result.Error != nil {
		return 0
	}
	return count
}

func (db *TransactionConnection) TotalTransactionByUserID(idUser uint64) int64 {
	var count int64
	result := db.connection.Model(&entity.Transaction{}).Where("id_user = ? && status = ?", idUser, 1).Count(&count)
	if result.Error != nil {
		return 0
	}
	return count
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionConnection{
		connection: db,
	}
}

func (db *TransactionConnection) InsertTransaction(Transaction *entity.Transaction) entity.Transaction {
	db.connection.Save(Transaction)
	return *Transaction
}

func (db *TransactionConnection) StorePaymentToken(transactionID uint64, paymentToken string, virtualAcc string) error {
	paymentTokenRecord := entity.PaymentToken{
		TransactionID: transactionID,
		PaymentToken:  paymentToken,
		VirtualAcc:    virtualAcc,
	}
	result := db.connection.Create(&paymentTokenRecord)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *TransactionConnection) All(page int, pageSize int) ([]entity.Transaction, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	var transactions []entity.Transaction
	offset := (page - 1) * pageSize

	result := db.connection.Where("status = ?", 1).Offset(offset).Limit(pageSize).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func (db *TransactionConnection) FindTransactionByIDUser(idUser uint64, page int, pageSize int) ([]entity.Transaction, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	var transactions []entity.Transaction
	offset := (page - 1) * pageSize

	result := db.connection.Where("id_user = ? && status = ?", idUser, 1).Offset(offset).Limit(pageSize).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func (db *TransactionConnection) UpdateTransaction(Transaction entity.Transaction) entity.Transaction {
	db.connection.Save(&Transaction)
	return Transaction
}

func (db *TransactionConnection) FindTransactionByID(id uint64) *entity.Transaction {
	var Transaction entity.Transaction
	result := db.connection.Where("id = ? ", id).Take(&Transaction)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}

	return &Transaction
}

func (db *TransactionConnection) UpdateTransactionStatus(id uint64, newStatus uint64) error {
	var trx entity.Transaction
	result := db.connection.First(&trx, id)
	if result.Error != nil {
		return result.Error
	}

	trx.Status = newStatus

	result = db.connection.Save(&trx)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
