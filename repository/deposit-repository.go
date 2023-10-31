package repository

import (
	"errors"

	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/helper"

	"gorm.io/gorm"
)

type DepositRepository interface {
	InsertDeposit(brg *entity.Deposit) entity.Deposit
	All(page int, pageSize int) ([]entity.Deposit, error)
	UpdateDeposit(plg entity.Deposit) entity.Deposit
	FindDepositByID(id string) entity.Deposit
	FindDepositByIDUser(id uint64, page int, pageSize int) ([]entity.Deposit, error)
	TotalDeposit() int64
	TotalDepositByUserID(idUser uint64) int64
	StorePaymentToken(depositID string, paymentToken string, virtualAcc string, callbackUrl string) error
	UpdateDepositStatus(id string, newStatus uint64) error
	FindPaymentInfoById(id string) *entity.PaymentToken
	SearchByDateAll(dateStart int64, dateEnd int64) ([]entity.Deposit, error)
	SearchByDateIDUser(idUser uint64, dateStart int64, dateEnd int64) ([]entity.Deposit, error)
	TotalDepositByDate(dateStart int64, dateEnd int64) int64
	TotalDepositByDateIdUser(idUser uint64, dateStart int64, dateEnd int64) int64
}

type DepositConnection struct {
	connection *gorm.DB
}

func (db *DepositConnection) TotalDeposit() int64 {
	var count int64
	result := db.connection.Model(&entity.Deposit{}).Where("status != ?", 1).Count(&count)
	if result.Error != nil {
		return 0
	}
	return count
}

func (db *DepositConnection) TotalDepositByUserID(idUser uint64) int64 {
	var count int64
	result := db.connection.Model(&entity.Deposit{}).Where("id_user = ? && status = ?", idUser, 5).Count(&count)
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
	Deposit.Date = helper.GetCurrentTimeInLocation()
	db.connection.Save(Deposit)
	return *Deposit
}

func (db *DepositConnection) StorePaymentToken(depositID string, paymentToken string, virtualAcc string, callbackUrl string) error {
	paymentTokenRecord := entity.PaymentToken{
		DepositID:    depositID,
		PaymentToken: paymentToken,
		VirtualAcc:   virtualAcc,
		CallbackUrl:  callbackUrl,
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

	result := db.connection.Where("status != ?", 1).Offset(offset).Limit(pageSize).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func (db *DepositConnection) SearchByDateAll(dateStart int64, dateEnd int64) ([]entity.Deposit, error) {
	var deposits []entity.Deposit

	condition := "date >= ? AND date <= ? AND status != ?"

	result := db.connection.Where(condition, dateStart, dateEnd, 1).Find(&deposits)
	if result.Error != nil {
		return nil, result.Error
	}

	return deposits, nil
}

func (db *DepositConnection) SearchByDateIDUser(idUser uint64, dateStart int64, dateEnd int64) ([]entity.Deposit, error) {
	var deposits []entity.Deposit

	condition := "id_user = ? AND date >= ? AND date <= ? AND status != ?"

	result := db.connection.Where(condition, idUser, dateStart, dateEnd, 1).Find(&deposits)
	if result.Error != nil {
		return nil, result.Error
	}

	return deposits, nil
}

func (db *DepositConnection) TotalDepositByDate(dateStart int64, dateEnd int64) int64 {
	var count int64

	condition := "date >= ? AND date <= ? AND status != ?"

	result := db.connection.Model(&entity.Deposit{}).Where(condition, dateStart, dateEnd, 1).Count(&count)
	if result.Error != nil {
		return 0
	}

	return count
}

func (db *DepositConnection) TotalDepositByDateIdUser(idUser uint64, dateStart int64, dateEnd int64) int64 {
	var count int64

	condition := "id_user = ? AND date >= ? AND date <= ? AND status != ?"

	result := db.connection.Model(&entity.Deposit{}).Where(condition, idUser, dateStart, dateEnd, 1).Count(&count)
	if result.Error != nil {
		return 0
	}

	return count
}

func (db *DepositConnection) FindDepositByIDUser(idUser uint64, page int, pageSize int) ([]entity.Deposit, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	var transactions []entity.Deposit
	offset := (page - 1) * pageSize

	result := db.connection.Where("id_user = ? && status != ?", idUser, 1).Offset(offset).Limit(pageSize).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func (db *DepositConnection) UpdateDeposit(Deposit entity.Deposit) entity.Deposit {
	db.connection.Save(&Deposit)
	return Deposit
}

func (db *DepositConnection) FindDepositByID(id string) entity.Deposit {
	var Deposit entity.Deposit
	result := db.connection.Where("id = ? ", id).Take(&Deposit)
	if result.Error != nil || result.RowsAffected == 0 {
		return Deposit
	}

	return Deposit
}

func (db *DepositConnection) FindPaymentInfoById(id string) *entity.PaymentToken {
	var payment entity.PaymentToken
	result := db.connection.Where("deposit_id = ? ", id).Take(&payment)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}

	return &payment
}

func (db *DepositConnection) UpdateDepositStatus(id string, newStatus uint64) error {
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
