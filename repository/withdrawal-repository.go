package repository

import (
	"errors"

	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/helper"

	"gorm.io/gorm"
)

type WithdrawalRepository interface {
	InsertWithdrawal(brg *entity.Withdrawal) entity.Withdrawal
	All(page int, pageSize int) ([]entity.Withdrawal, error)
	UpdateWithdrawal(plg entity.Withdrawal) entity.Withdrawal
	FindWithdrawalByID(id uint64) *entity.Withdrawal
	FindWithdrawalByIDUser(id uint64, page int, pageSize int) ([]entity.Withdrawal, error)
	TotalWithdrawal() int64
	TotalWithdrawalByUserID(idUser uint64) int64
	UpdateWithdrawalStatus(id uint64, newStatus uint64) error
}

type WithdrawalConnection struct {
	connection *gorm.DB
}

func (db *WithdrawalConnection) TotalWithdrawal() int64 {
	var count int64
	result := db.connection.Model(&entity.Withdrawal{}).Where("status = ?", 1).Count(&count)
	if result.Error != nil {
		return 0
	}
	return count
}

func (db *WithdrawalConnection) TotalWithdrawalByUserID(idUser uint64) int64 {
	var count int64
	result := db.connection.Model(&entity.Withdrawal{}).Where("id_user = ? && status = ?", idUser, 1).Count(&count)
	if result.Error != nil {
		return 0
	}
	return count
}

func NewWithdrawalRepository(db *gorm.DB) WithdrawalRepository {
	return &WithdrawalConnection{
		connection: db,
	}
}

func (db *WithdrawalConnection) InsertWithdrawal(Withdrawal *entity.Withdrawal) entity.Withdrawal {
	Withdrawal.Date = helper.GetCurrentTimeInLocation()
	db.connection.Save(Withdrawal)
	return *Withdrawal
}

func (db *WithdrawalConnection) All(page int, pageSize int) ([]entity.Withdrawal, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	var transactions []entity.Withdrawal
	offset := (page - 1) * pageSize

	result := db.connection.Where("status = ?", 1).Offset(offset).Limit(pageSize).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func (db *WithdrawalConnection) FindWithdrawalByIDUser(idUser uint64, page int, pageSize int) ([]entity.Withdrawal, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	var transactions []entity.Withdrawal
	offset := (page - 1) * pageSize

	result := db.connection.Where("id_user = ? && status = ?", idUser, 1).Offset(offset).Limit(pageSize).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func (db *WithdrawalConnection) UpdateWithdrawal(Withdrawal entity.Withdrawal) entity.Withdrawal {
	db.connection.Save(&Withdrawal)
	return Withdrawal
}

func (db *WithdrawalConnection) FindWithdrawalByID(id uint64) *entity.Withdrawal {
	var Withdrawal entity.Withdrawal
	result := db.connection.Where("id = ? ", id).Take(&Withdrawal)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}

	return &Withdrawal
}

func (db *WithdrawalConnection) UpdateWithdrawalStatus(id uint64, newStatus uint64) error {
	var trx entity.Withdrawal
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
