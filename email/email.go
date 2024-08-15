package email

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/types"
)

func SendWelcomeEmail(user types.User) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	from := "server@tb-books.com"
	to := []string{user.Email}
	subject := "Welcome, " + user.Username
	body := "Welcome to TB Books, " + user.Username + "!"

	msg := []byte("To: " + to[0] + "\r\n" + "Subject: " + subject + "\r\n" + "\r\n" + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, nil, from, to, msg)
	if err != nil {
		fmt.Println("Error sending welcome email:", err)
		return
	}

	fmt.Println("Welcome mail sent successfully to " + user.Email)
}

func SendEmailVerification(user types.User, token string) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	from := "server@tb-books.com"
	to := []string{user.Email}
	subject := "Verify your email"
	body := "Welcome to TB Books, " + user.Username + "! Please verify your email by clicking the following link: http://" + os.Getenv("BASE_URL_API") + "/api/users/verify-email/" + token

	msg := []byte("To: " + to[0] + "\r\n" + "Subject: " + subject + "\r\n" + "\r\n" + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, nil, from, to, msg)
	if err != nil {
		fmt.Println("Error sending verification email:", err)
		return
	}

	fmt.Println("Verification email sent successfully to " + user.Email)
}
