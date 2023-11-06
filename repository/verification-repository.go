package repository

import (
	"time"

	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type VerificationRepository interface {
	InsertVerification(email string, verificationKey string) error
	ValidateVerification(verificationKey string) bool
}

type redisConnection struct {
	connection   *redis.Client
	connectionDB *gorm.DB
}

func NewVerificationRepository(db *redis.Client, sqlDB *gorm.DB) VerificationRepository {
	return &redisConnection{connection: db, connectionDB: sqlDB}
}

func (db *redisConnection) InsertVerification(email string, verificationKey string) error {
	statusCMD := db.connection.Set(verificationKey, email, time.Minute*3)
	if statusCMD.Err() != nil {
		logrus.Error(statusCMD.Err())
		return statusCMD.Err()
	}

	res, err := statusCMD.Result()
	if err != nil {
		logrus.Error(err.Error())
	}

	logrus.Info("OTP Inserted to Redis ", res)

	return nil
}

func (db *redisConnection) ValidateVerification(verificationKey string) bool {
	email, statusCMD := db.connection.Get(verificationKey).Result()
	if statusCMD != nil {
		logrus.Error(statusCMD.Error())
		return false
	}

	var user entity.User
	if err := db.connectionDB.Where("email = ?", email).First(&user).Error; err != nil {
		return true
	}

	user.IsVerified = true

	if err := db.connectionDB.Save(&user).Error; err != nil {
		return true
	}

	_, err := db.connection.Del(verificationKey).Result()
	if err != nil {
		return true
	}
	return true
}
