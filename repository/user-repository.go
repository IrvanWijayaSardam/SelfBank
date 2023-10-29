package repository

import (
	"log"

	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/helper"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	All() []entity.User
	InsertUser(user entity.User) entity.User
	UpdateUser(user entity.User) entity.User
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

func (db *userConnection) All() []entity.User {
	var users []entity.User
	db.connection.Where("status = ?", 1).Find(&users)
	return users
}

func (db *userConnection) InsertUser(user entity.User) entity.User {
	user.Password = hashAndSalt([]byte(user.Password))
	user.AccountNumber = helper.GenerateRandomAccountNumber()
	db.connection.Save(&user)
	return user
}

func (db *userConnection) UpdateUser(user entity.User) entity.User {
	if user.Password != "" {
		user.Password = hashAndSalt([]byte(user.Password))
	} else {
		var tempUser entity.User
		db.connection.Find(&tempUser, user.ID)
		user.Password = tempUser.Password
	}

	db.connection.Save(&user)
	return user
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

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}
