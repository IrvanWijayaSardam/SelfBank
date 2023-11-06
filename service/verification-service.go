package service

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"

	"github.com/IrvanWijayaSardam/SelfBank/helper"
	"github.com/IrvanWijayaSardam/SelfBank/repository"
)

type VerificationService interface {
	SendVerificationEmail(email string) error
	VerifyOtp(otp string) bool
}

type verificationService struct {
	VerificationRepository repository.VerificationRepository
}

func NewVerificationService(verifRepo repository.VerificationRepository) VerificationService {
	return &verificationService{
		VerificationRepository: verifRepo,
	}
}

func (service *verificationService) SendVerificationEmail(email string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_MAIL")
	password := os.Getenv("SMTP_PASSWORD")
	var otp = strconv.FormatUint(helper.GenerateRandomAccountNumber(), 10)

	auth := smtp.PlainAuth("", from, password, host)
	smtpAddr := fmt.Sprintf("%s:%s", host, port)

	body := "From: " + from + "\n" +
		"To: " + email + "\n" +
		"Subject: " + "Email Verification" + "\n\n" +
		otp

	err := smtp.SendMail(smtpAddr, auth, from, []string{email}, []byte(body))
	if err != nil {
		return err
	}

	err = service.VerificationRepository.InsertVerification(email, otp)
	if err != nil {
		return err
	}

	return nil
}

func (service *verificationService) VerifyOtp(token string) bool {
	return service.VerificationRepository.ValidateVerification(token)
}
