package email

import (
	"fmt"
	"net/smtp"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/types"
)

func SendEmail(user types.User) {
	smtpHost := "localhost"
	smtpPort := "1025"

	from := "server@tbbooks.com"
	to := []string{user.Email}
	subject := "Welcome, " + user.Username
	body := "Welcome to TB Books, " + user.Username + "!"

	msg := []byte("To: " + to[0] + "\r\n" + "Subject: " + subject + "\r\n" + "\r\n" + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, nil, from, to, msg)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully")
}

func SendEmailVerification(user types.User, token string) {
	smtpHost := "localhost"
	smtpPort := "1025"

	from := "server@tbbooks.com"
	to := []string{user.Email}
	subject := "Verify your email"
	body := "Welcome to TB Books, " + user.Username + "! Please verify your email by clicking the following link: http://localhost:2337/api/users/verify-email/" + token

	msg := []byte("To: " + to[0] + "\r\n" + "Subject: " + subject + "\r\n" + "\r\n" + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, nil, from, to, msg)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully")
}
