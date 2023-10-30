package repository

import (
	"errors"

	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/helper"

	"gorm.io/gorm"
)

type UserRepository interface {
	All(page int, pageSize int) ([]entity.User, error)
	InsertUser(user entity.User) entity.User
	UpdateUser(user entity.User) entity.User
	DeleteUser(idUser uint64) bool
	VerifyCredential(email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	FindByEmail(email string) entity.User
	ProfileUser(userId uint64) entity.User
	TotalDepositByUserID(userId uint64) int64
	TotalWithdrawalByUserID(userid uint64) int64
	TotalTransactionInByAccountNumber(accountNumber string) int64
	TotalTransactionFromByAccountNumber(accountNumber string) int64
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) All(page int, pageSize int) ([]entity.User, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, errors.New("Invalid page or pageSize values")
	}

	var transactions []entity.User
	offset := (page - 1) * pageSize

	result := db.connection.Where("status = ?", 1).Offset(offset).Limit(pageSize).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func (db *userConnection) InsertUser(user entity.User) entity.User {
	user.Password = helper.HashAndSalt([]byte(user.Password))
	user.AccountNumber = helper.GenerateRandomAccountNumber()
	db.connection.Save(&user)
	return user
}

func (db *userConnection) UpdateUser(user entity.User) entity.User {
	db.connection.Save(&user)
	return user
}

func (db *userConnection) DeleteUser(idUser uint64) bool {
	var user entity.User
	db.connection.Where("id = ?", idUser).Take(&user)
	user.Status = 2
	result := db.connection.Save(&user)

	if result.RowsAffected == 1 {
		return true
	}
	return false
}

func (db *userConnection) VerifyCredential(email string, password string) interface{} {
	var user entity.User
	res := db.connection.Where("email =?", email).Take(&user)
	if res.Error == nil {
		return user
	}
	return nil
}

func (db *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user entity.User
	return db.connection.Where("email = ?", email).Take(&user)
}

func (db *userConnection) FindByEmail(email string) entity.User {
	var user entity.User
	db.connection.Where("email LIKE ? AND status = ?", "%"+email+"%", 1).Take(&user)
	return user
}

func (db *userConnection) ProfileUser(userID uint64) entity.User {
	var user entity.User
	db.connection.Find(&user, userID)
	return user
}

func (db *userConnection) TotalDepositByUserID(idUser uint64) int64 {
	var totalAmount int64
	result := db.connection.Model(&entity.Deposit{}).Select("SUM(amount)").Where("id_user = ? && status = ?", idUser, 5).Scan(&totalAmount)
	if result.Error != nil {
		return 0
	}
	return totalAmount
}

func (db *userConnection) TotalWithdrawalByUserID(idUser uint64) int64 {
	var totalAmount int64
	result := db.connection.Model(&entity.Withdrawal{}).Select("SUM(amount)").Where("id_user = ? && status = ?", idUser, 1).Scan(&totalAmount)
	if result.Error != nil {
		return 0
	}
	return totalAmount
}

func (db *userConnection) TotalTransactionInByAccountNumber(accountNumber string) int64 {
	var totalAmount int64
	result := db.connection.Model(&entity.Transaction{}).Select("SUM(amount)").Where("transaction_to = ? && status = ?", accountNumber, 1).Scan(&totalAmount)
	if result.Error != nil {
		return 0
	}
	return totalAmount
}

func (db *userConnection) TotalTransactionFromByAccountNumber(accountNumber string) int64 {
	var totalAmount int64
	result := db.connection.Model(&entity.Transaction{}).Select("SUM(amount)").Where("transaction_from = ? && status = ?", accountNumber, 1).Scan(&totalAmount)
	if result.Error != nil {
		return 0
	}
	return totalAmount
}
