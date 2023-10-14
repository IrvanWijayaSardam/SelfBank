package repository

import (
	"errors"

	"github.com/IrvanWijayaSardam/SelfBank/entity"

	"gorm.io/gorm"
)

type DepositRepository interface {
	InsertDeposit(brg *entity.Deposit) entity.Deposit
	All(page int, pageSize int) ([]entity.Deposit, error)
	UpdateDeposit(plg entity.Deposit) entity.Deposit
	FindDepositByID(id uint64) *entity.Deposit
	FindDepositByIDUser(id uint64, page int, pageSize int) ([]entity.Deposit, error)
	TotalDeposit() int64
	TotalDepositByUserID(idUser uint64) int64
	StorePaymentToken(transactionID uint64, paymentToken string, virtualAcc string) error
	UpdateDepositStatus(id uint64, newStatus uint64) error
}

type DepositConnection struct {
	connection *gorm.DB
}

func (db *DepositConnection) TotalDeposit() int64 {
	var count int64
	result := db.connection.Model(&entity.Deposit{}).Where("status = ?", 1).Count(&count)
	if result.Error != nil {
		return 0
	}
	return count
}

func (db *DepositConnection) TotalDepositByUserID(idUser uint64) int64 {
	var count int64
	result := db.connection.Model(&entity.Deposit{}).Where("id_user = ? && status = ?", idUser, 1).Count(&count)
	if result.Error != nil {
		return 0
	}
	return count
}

func NewDepositRepository(db *gorm.DB) DepositRepository {
	return &DepositConnection{
		connection: db,
	}
}

func (db *DepositConnection) InsertDeposit(Deposit *entity.Deposit) entity.Deposit {
	db.connection.Save(Deposit)
	return *Deposit
}

func (db *DepositConnection) StorePaymentToken(transactionID uint64, paymentToken string, virtualAcc string) error {
	paymentTokenRecord := entity.PaymentToken{
		DepositID:    transactionID,
		PaymentToken: paymentToken,
		VirtualAcc:   virtualAcc,
	}
	result := db.connection.Create(&paymentTokenRecord)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DepositConnection) All(page int, pageSize int) ([]entity.Deposit, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	var transactions []entity.Deposit
	offset := (page - 1) * pageSize

	result := db.connection.Where("status = ?", 1).Offset(offset).Limit(pageSize).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func (db *DepositConnection) FindDepositByIDUser(idUser uint64, page int, pageSize int) ([]entity.Deposit, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	var transactions []entity.Deposit
	offset := (page - 1) * pageSize

	result := db.connection.Where("id_user = ? && status = ?", idUser, 1).Offset(offset).Limit(pageSize).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func (db *DepositConnection) UpdateDeposit(Deposit entity.Deposit) entity.Deposit {
	db.connection.Save(&Deposit)
	return Deposit
}

func (db *DepositConnection) FindDepositByID(id uint64) *entity.Deposit {
	var Deposit entity.Deposit
	result := db.connection.Where("id = ? ", id).Take(&Deposit)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}

	return &Deposit
}

func (db *DepositConnection) UpdateDepositStatus(id uint64, newStatus uint64) error {
	var trx entity.Deposit
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
