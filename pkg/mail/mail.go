package mail

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/smtp"
	"regexp"
	"time"

	"github.com/CollabTED/CollabTed-Backend/config"
	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/CollabTED/CollabTed-Backend/pkg/redis"
	r "github.com/redis/go-redis/v9"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

type EmailVerifier struct {
	client *r.Client
}

func NewVerifier() *EmailVerifier {
	return &EmailVerifier{
		client: redis.GetClient(),
	}
}

func (v *EmailVerifier) GenerateOTP() string {
	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if n != 6 {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func (v *EmailVerifier) SendVerfication(userID string, to []string) error {
	smtpHost := config.EMAIL_HOST
	smtpPort := config.EMAIL_PORT
	otp := v.GenerateOTP()
	message := []byte(fmt.Sprintf("Verification code is %s", otp))

	// Log OTP being set
	logger.LogInfo().Msg(fmt.Sprintf("Storing OTP: %s for key: %s", otp, userID))
	v.client.Set(context.Background(), userID, otp, time.Hour*1)

	auth := smtp.PlainAuth("", config.EMAIL, config.EMAIL_PASSWORD, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, config.EMAIL, to, message)
	if err != nil {
		return err
	}
	return nil
}

func (v *EmailVerifier) Verify(userID string, otp string) error {
	logger.LogDebug().Msg(fmt.Sprintf("Verifying OTP: %s for userID: %s", otp, userID))

	userOTP := v.client.Get(context.Background(), userID).Val()
	logger.LogDebug().Msg(fmt.Sprintf("Retrieved OTP from cache: %s", userOTP))

	if userOTP == "" {
		return errors.New("verification failed: OTP not found")
	}

	if userOTP != otp {
		return errors.New("verification failed: OTP does not match")
	}
	return nil
}

func IsValidEmail(email string) bool {
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$` // this should be updated later
	re := regexp.MustCompile(emailPattern)
	return re.MatchString(email)
}
